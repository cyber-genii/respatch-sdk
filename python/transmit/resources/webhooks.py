from typing import Any, Dict, List, Optional


class Webhooks:
    def __init__(self, client):
        self._client = client

    def list(self) -> List[Dict[str, Any]]:
        """List your webhooks"""
        return self._client.get("/api/v1/developers/webhooks")

    def create(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """Create a webhook"""
        return self._client.post("/api/v1/developers/webhooks", json=data)

    def retrieve(self, webhook_id: str) -> Dict[str, Any]:
        """Get webhook details"""
        return self._client.get(f"/api/v1/developers/webhooks/{webhook_id}")

    def update(self, webhook_id: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Update webhook"""
        return self._client.put(f"/api/v1/developers/webhooks/{webhook_id}", json=data)

    def delete(self, webhook_id: str) -> None:
        """Delete webhook"""
        self._client.delete(f"/api/v1/developers/webhooks/{webhook_id}")

    def get_events(self, webhook_id: str) -> List[Dict[str, Any]]:
        """Get webhook events"""
        return self._client.get(f"/api/v1/developers/webhooks/{webhook_id}/events")

    def dlq(self, params: Optional[Dict[str, Any]] = None) -> Any:
        """List dead-letter / failed webhook events"""
        return self._client.get("/api/v1/developers/webhook-events", params=params)

    def replay_dlq(self, event_id: str) -> Any:
        """Retry / replay a failed webhook event from the DLQ"""
        return self._client.post(f"/api/v1/developers/webhook-events/{event_id}/retry")
