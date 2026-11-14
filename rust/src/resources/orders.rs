use crate::client::{Client, RespatchError};
use crate::types::{CreateDeliveryOrderRequest, GetQuoteRequest};
use serde_json::Value;

pub struct Orders<'a> {
    client: &'a Client,
}

impl<'a> Orders<'a> {
    pub fn new(client: &'a Client) -> Self {
        Self { client }
    }

    pub async fn quote(&self, req: &GetQuoteRequest) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/delivery-orders/quote", req).await
    }

    /// Deprecated — use `quote`.
    pub async fn calculate_fare(&self, req: &GetQuoteRequest) -> Result<Value, RespatchError> {
        self.client
            .post("/api/v1/delivery-orders/calculate-fare", req)
            .await
    }

    pub async fn create(&self, req: &CreateDeliveryOrderRequest) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/delivery-orders", req).await
    }

    /// Create an order from a prior quote (single HTTP call).
    pub async fn book_from_quote(
        &self,
        req: &CreateDeliveryOrderRequest,
    ) -> Result<Value, RespatchError> {
        self.create(req).await
    }

    pub async fn list(&self) -> Result<Vec<Value>, RespatchError> {
        self.client.get("/api/v1/delivery-orders").await
    }

    pub async fn retrieve(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}", id);
        self.client.get(&path).await
    }

    pub async fn track(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}/track", id);
        self.client.get(&path).await
    }

    pub async fn cancel(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}/cancel", id);
        self.client.post(&path, &serde_json::json!({})).await
    }

    pub async fn generate_auth_code(&self, id: &str, action: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}/auth-code/generate", id);
        self.client
            .post(&path, &serde_json::json!({ "action": action }))
            .await
    }

    pub async fn get_active_auth_code(&self, id: &str, action: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}/auth-code/active?action={}", id, action);
        self.client.get(&path).await
    }
}
