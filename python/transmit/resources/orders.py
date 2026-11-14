from typing import Any, Dict, List, Optional, TypedDict


class QuoteLocationInput(TypedDict, total=False):
    address: str
    lat: float
    lng: float


class PackageDimensions(TypedDict):
    length: float
    width: float
    height: float


class QuotePackageInput(TypedDict, total=False):
    weight_kg: float
    dimensions_cm: PackageDimensions
    quantity: int


class GetQuoteRequest(TypedDict, total=False):
    pickup: QuoteLocationInput
    dropoff: QuoteLocationInput
    delivery_type: str
    vehicle_type: str
    packages: List[QuotePackageInput]
    scheduled_pickup_time: str


class OrderContactLocation(TypedDict, total=False):
    contact_name: str
    contact_email: str
    contact_phone: str
    contact_phone_secondary: str


class OrderPackageMetaInput(TypedDict, total=False):
    package_id: str
    description: str
    value: float
    fragile: bool
    supplier_info: Dict[str, str]


class CreateDeliveryOrderRequest(TypedDict, total=False):
    quote_id: str
    pickup: OrderContactLocation
    dropoff: OrderContactLocation
    packages: List[OrderPackageMetaInput]
    payment_method: str
    webhook_url: str
    external_reference: str


# Deprecated aliases
CalculateOrderFareRequest = GetQuoteRequest


class Orders:
    def __init__(self, client):
        self._client = client

    def quote(self, data: GetQuoteRequest) -> Dict[str, Any]:
        return self._client.post("/api/v1/delivery-orders/quote", json=data)

    def calculate_fare(self, data: GetQuoteRequest) -> Dict[str, Any]:
        """Deprecated — use quote()."""
        return self._client.post("/api/v1/delivery-orders/calculate-fare", json=data)

    def create(self, data: CreateDeliveryOrderRequest) -> Dict[str, Any]:
        return self._client.post("/api/v1/delivery-orders", json=data)

    def book_from_quote(
        self,
        quote_id: str,
        pickup: OrderContactLocation,
        dropoff: OrderContactLocation,
        packages: List[OrderPackageMetaInput],
        **kwargs: Any,
    ) -> Dict[str, Any]:
        """Create an order from a prior quote (single HTTP call)."""
        return self.create(
            {
                "quote_id": quote_id,
                "pickup": pickup,
                "dropoff": dropoff,
                "packages": packages,
                **kwargs,
            }
        )

    def list(self, params: Optional[Dict[str, Any]] = None) -> List[Dict[str, Any]]:
        return self._client.get("/api/v1/delivery-orders", params=params)

    def retrieve(self, order_id: str) -> Dict[str, Any]:
        return self._client.get(f"/api/v1/delivery-orders/{order_id}")

    def track(self, order_id: str) -> Dict[str, Any]:
        return self._client.get(f"/api/v1/delivery-orders/{order_id}/track")

    def cancel(self, order_id: str) -> None:
        self._client.post(f"/api/v1/delivery-orders/{order_id}/cancel")

    def generate_auth_code(self, order_id: str, action: str = "pickup") -> Dict[str, Any]:
        return self._client.post(
            f"/api/v1/delivery-orders/{order_id}/auth-code/generate",
            json={"action": action},
        )

    def get_active_auth_code(self, order_id: str, action: str = "pickup") -> Dict[str, Any]:
        return self._client.get(
            f"/api/v1/delivery-orders/{order_id}/auth-code/active",
            params={"action": action},
        )
