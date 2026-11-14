use crate::client::{Client, RespatchError};
use serde_json::Value;

pub struct Webhooks<'a> {
    client: &'a Client,
}

impl<'a> Webhooks<'a> {
    pub fn new(client: &'a Client) -> Self {
        Self { client }
    }

    pub async fn list(&self) -> Result<Vec<Value>, RespatchError> {
        self.client.get("/api/v1/developers/webhooks").await
    }

    pub async fn create(&self, body: &Value) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/developers/webhooks", body).await
    }

    pub async fn retrieve(&self, id: &str) -> Result<Value, RespatchError> {
        self.client
            .get(&format!("/api/v1/developers/webhooks/{}", id))
            .await
    }

    pub async fn delete(&self, id: &str) -> Result<Value, RespatchError> {
        self.client
            .delete(&format!("/api/v1/developers/webhooks/{}", id))
            .await
    }

    pub async fn dlq(
        &self,
        status: Option<&str>,
        page: Option<u64>,
        per_page: Option<u64>,
    ) -> Result<Value, RespatchError> {
        let mut params = Vec::new();
        if let Some(s) = status {
            params.push(format!("status={}", s));
        }
        if let Some(p) = page {
            params.push(format!("page={}", p));
        }
        if let Some(pp) = per_page {
            params.push(format!("per_page={}", pp));
        }
        let path = if params.is_empty() {
            "/api/v1/developers/webhook-events".to_string()
        } else {
            format!("/api/v1/developers/webhook-events?{}", params.join("&"))
        };
        self.client.get(&path).await
    }

    pub async fn replay_dlq(&self, event_id: &str) -> Result<Value, RespatchError> {
        self.client
            .post(
                &format!("/api/v1/developers/webhook-events/{}/retry", event_id),
                &serde_json::json!({}),
            )
            .await
    }
}
