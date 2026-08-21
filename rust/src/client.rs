use crate::response::{decode_data, unwrap_value};
use reqwest::{Client as ReqwestClient, Response, header};
use serde::{de::DeserializeOwned, Serialize};
use serde_json::Value;
use std::fmt;

#[derive(Debug)]
pub enum TransmitError {
    Request(reqwest::Error),
    Api { status: u16, message: String },
    Json(serde_json::Error),
}

impl fmt::Display for TransmitError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            TransmitError::Request(e) => write!(f, "Request error: {}", e),
            TransmitError::Api { status, message } => write!(f, "API error [{}]: {}", status, message),
            TransmitError::Json(e) => write!(f, "JSON error: {}", e),
        }
    }
}

impl std::error::Error for TransmitError {}

impl From<reqwest::Error> for TransmitError {
    fn from(err: reqwest::Error) -> Self {
        TransmitError::Request(err)
    }
}

impl From<serde_json::Error> for TransmitError {
    fn from(err: serde_json::Error) -> Self {
        TransmitError::Json(err)
    }
}

#[derive(Clone)]
pub struct Client {
    pub(crate) http: ReqwestClient,
    pub(crate) base_url: String,
}

impl Client {
    pub fn new(api_key: &str, environment: Option<&str>, base_url_override: Option<&str>) -> Result<Self, TransmitError> {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            header::AUTHORIZATION,
            header::HeaderValue::from_str(&format!("Bearer {}", api_key))
                .map_err(|_| TransmitError::Api { status: 0, message: "Invalid API key format".into() })?,
        );
        headers.insert(header::CONTENT_TYPE, header::HeaderValue::from_static("application/json"));
        headers.insert(header::USER_AGENT, header::HeaderValue::from_static("Respatch-Rust-SDK/1.0.0"));

        let http_client = ReqwestClient::builder()
            .default_headers(headers)
            .build()?;

        let base_url = if let Some(custom_url) = base_url_override {
            custom_url.to_string()
        } else if environment == Some("sandbox") {
            // respatch-sandbox frontier
            "https://sandbox-api.respatch.com".to_string()
        } else {
            // respatch-api frontier
            "https://api.respatch.com".to_string()
        };

        Ok(Self {
            http: http_client,
            base_url,
        })
    }

    pub(crate) async fn handle_response<T: DeserializeOwned>(response: Response) -> Result<T, TransmitError> {
        let status = response.status();
        if !status.is_success() {
            let error_text = response.text().await.unwrap_or_else(|_| "Unknown error".to_string());
            let message = match serde_json::from_str::<Value>(&error_text) {
                Ok(json) => {
                    if let Some(msg) = json.get("message").and_then(|m| m.as_str()) {
                        msg.to_string()
                    } else if let Some(err) = json.get("error").and_then(|e| e.as_str()) {
                        err.to_string()
                    } else {
                        error_text
                    }
                }
                Err(_) => error_text,
            };
            return Err(TransmitError::Api {
                status: status.as_u16(),
                message,
            });
        }
        
        let raw: Value = response.json().await?;
        decode_data(unwrap_value(raw)).map_err(TransmitError::from)
    }

    pub(crate) async fn get<T: DeserializeOwned>(&self, path: &str) -> Result<T, TransmitError> {
        let url = format!("{}{}", self.base_url, path);
        let resp = self.http.get(&url).send().await?;
        Self::handle_response(resp).await
    }

    pub(crate) async fn post<T: DeserializeOwned, B: Serialize>(&self, path: &str, body: &B) -> Result<T, TransmitError> {
        let url = format!("{}{}", self.base_url, path);
        let resp = self.http.post(&url).json(body).send().await?;
        Self::handle_response(resp).await
    }

    pub(crate) async fn put<T: DeserializeOwned, B: Serialize>(&self, path: &str, body: &B) -> Result<T, TransmitError> {
        let url = format!("{}{}", self.base_url, path);
        let resp = self.http.put(&url).json(body).send().await?;
        Self::handle_response(resp).await
    }

    pub(crate) async fn delete<T: DeserializeOwned>(&self, path: &str) -> Result<T, TransmitError> {
        let url = format!("{}{}", self.base_url, path);
        let resp = self.http.delete(&url).send().await?;
        Self::handle_response(resp).await
    }
}
