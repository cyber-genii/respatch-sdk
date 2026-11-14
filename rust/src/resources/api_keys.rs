use crate::client::{Client, RespatchError};
use crate::types::ApiKeyCreateRequest;
use serde_json::Value;

pub struct ApiKeys<'a> {
    client: &'a Client,
}

impl<'a> ApiKeys<'a> {
    pub fn new(client: &'a Client) -> Self {
        Self { client }
    }

    pub async fn list(&self) -> Result<Vec<Value>, RespatchError> {
        self.client.get("/api/v1/developers/api-keys").await
    }

    pub async fn create(&self, req: &ApiKeyCreateRequest) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/developers/api-keys", req).await
    }

    pub async fn retrieve(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/developers/api-keys/{}", id);
        self.client.get(&path).await
    }

    pub async fn revoke(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/developers/api-keys/{}/revoke", id);
        self.client.post(&path, &serde_json::json!({})).await
    }

    pub async fn delete(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/developers/api-keys/{}", id);
        self.client.delete(&path).await
    }

    pub async fn get_usage(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/developers/api-keys/{}/usage", id);
        self.client.get(&path).await
    }
}
