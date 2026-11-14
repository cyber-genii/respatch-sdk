use crate::client::{Client, RespatchError};
use crate::types::{CalculateOrderFareRequest, CreateDeliveryOrderRequest};
use serde_json::Value;

pub struct Deliveries<'a> {
    client: &'a Client,
}

impl<'a> Deliveries<'a> {
    pub fn new(client: &'a Client) -> Self {
        Self { client }
    }

    pub async fn create(&self, req: &CreateDeliveryOrderRequest) -> Result<Value, RespatchError> {
        self.client.post("/api/v1/delivery-orders", req).await
    }

    pub async fn list(&self) -> Result<Vec<Value>, RespatchError> {
        self.client.get("/api/v1/delivery-orders").await
    }

    pub async fn retrieve(&self, id: &str) -> Result<Value, RespatchError> {
        let path = format!("/api/v1/delivery-orders/{}", id);
        self.client.get(&path).await
    }

    pub async fn quote(&self, req: &CalculateOrderFareRequest) -> Result<Value, RespatchError> {
        self.client
            .post("/api/v1/delivery-orders/calculate-fare", req)
            .await
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
