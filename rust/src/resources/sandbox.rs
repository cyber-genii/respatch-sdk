use crate::client::{Client, RespatchError};
use serde_json::Value;

pub struct Sandbox<'a> {
    client: &'a Client,
}

impl<'a> Sandbox<'a> {
    pub fn new(client: &'a Client) -> Self {
        Self { client }
    }

    pub async fn create_data(&self, body: &Value) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/developers/sandbox", body).await
    }

    pub async fn list_data(&self) -> Result<Vec<Value>, RespatchError> {
        self.client.get("/api/v1/developers/sandbox").await
    }

    pub async fn clear_data(&self) -> Result<Value, RespatchError> {
        self.client.delete("/api/v1/developers/sandbox").await
    }
}
