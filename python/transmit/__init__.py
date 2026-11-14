from .client import RespatchAPIError, RespatchClient, TransmitAPIError, TransmitClient
from .resources.api_keys import ApiKeys
from .resources.deliveries import Deliveries
from .resources.orders import Orders
from .resources.sandbox import Sandbox
from .resources.vehicles import Vehicles
from .resources.webhooks import Webhooks


class Respatch:
    """Main Respatch SDK client."""

    def __init__(self, api_key: str = None, **kwargs):
        self._client = RespatchClient(api_key=api_key, **kwargs)
        self.orders = Orders(self._client)
        self.deliveries = Deliveries(self._client)
        self.api_keys = ApiKeys(self._client)
        self.webhooks = Webhooks(self._client)
        self.sandbox = Sandbox(self._client)
        self.vehicles = Vehicles(self._client)


Transmit = Respatch

from .webhook import (
    DEFAULT_WEBHOOK_TOLERANCE_SECONDS,
    RESPATCH_SIGNATURE_HEADER,
    get_webhook_signature_header,
    sign_webhook_header,
    sign_webhook_payload,
    verify_webhook_signature,
)

__all__ = [
    "Respatch",
    "Transmit",
    "RespatchClient",
    "TransmitClient",
    "RespatchAPIError",
    "TransmitAPIError",
    "verify_webhook_signature",
    "sign_webhook_payload",
    "sign_webhook_header",
    "get_webhook_signature_header",
    "RESPATCH_SIGNATURE_HEADER",
    "DEFAULT_WEBHOOK_TOLERANCE_SECONDS",
]
