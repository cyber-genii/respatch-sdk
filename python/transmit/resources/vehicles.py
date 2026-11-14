from typing import Any, Dict, List


class Vehicles:
    def __init__(self, client):
        self._client = client

    def get_types(self) -> List[Dict[str, Any]]:
        """Get available vehicle types"""
        return self._client.get("/api/v1/vehicles/types")

    def truck_specialties(self) -> List[Dict[str, Any]]:
        """Get available truck specialties and surcharge metadata"""
        return self._client.get("/api/v1/vehicles/truck-specialties")
