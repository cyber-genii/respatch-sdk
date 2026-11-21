from typing import Any, Dict, List


class Sandbox:
    def __init__(self, client):
        self._client = client

    def create_data(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Create sandbox test data"""
        return self._client.post("/api/v1/developers/sandbox", json=data)

    def list_data(self) -> List[Dict[str, Any]]:
        """List sandbox data"""
        return self._client.get("/api/v1/developers/sandbox")

    def clear_data(self) -> None:
        """Clear all sandbox data"""
        self._client.delete("/api/v1/developers/sandbox")

    def simulate_event(self, order_id: str, event: str) -> Dict[str, Any]:
        """Simulate a delivery lifecycle state transition in the sandbox"""
        return self._client.post(
            f"/api/v1/sandbox/orders/{order_id}/simulate/{event}",
            json={},
        )
