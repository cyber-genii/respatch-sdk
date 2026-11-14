from typing import Any, Dict, List


class Deliveries:
    """Deprecated — use ``Orders`` instead."""

    def __init__(self, client):
        self._client = client

    def create(self, data: Dict[str, Any]) -> Dict[str, Any]:
        return self._client.post("/api/v1/delivery-orders", json=data)

    def list(self, params: Dict[str, Any] = None) -> List[Dict[str, Any]]:
        return self._client.get("/api/v1/delivery-orders", params=params)

    def retrieve(self, delivery_id: str) -> Dict[str, Any]:
        return self._client.get(f"/api/v1/delivery-orders/{delivery_id}")

    def quote(self, data: Dict[str, Any]) -> Dict[str, Any]:
        return self._client.post("/api/v1/delivery-orders/calculate-fare", json=data)

    def generate_auth_code(self, delivery_id: str, action: str = "pickup") -> Dict[str, Any]:
        return self._client.post(
            f"/api/v1/delivery-orders/{delivery_id}/auth-code/generate",
            json={"action": action},
        )

    def get_active_auth_code(self, delivery_id: str, action: str = "pickup") -> Dict[str, Any]:
        return self._client.get(
            f"/api/v1/delivery-orders/{delivery_id}/auth-code/active",
            params={"action": action},
        )
