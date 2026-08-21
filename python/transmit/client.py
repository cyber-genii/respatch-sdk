import os
from typing import Any, Dict, Optional

import requests

from .response import unwrap_data

PRODUCTION_URL = "https://api.respatch.com"  # respatch-api frontier
SANDBOX_URL = "https://sandbox-api.respatch.com"  # respatch-sandbox frontier


class RespatchAPIError(Exception):
    """Exception raised for errors in the Respatch API."""

    def __init__(self, message, status_code=None, response=None):
        super().__init__(message)
        self.status_code = status_code
        self.response = response


# Backward compatibility
TransmitAPIError = RespatchAPIError


class RespatchClient:
    def __init__(
        self,
        api_key: Optional[str] = None,
        environment: str = "production",
        base_url: Optional[str] = None,
        timeout: int = 30,
    ):
        self.api_key = api_key or os.environ.get("RESPATCH_API_KEY") or os.environ.get("TRANSMIT_API_KEY")
        if not self.api_key:
            raise ValueError("API Key is required (pass api_key or set RESPATCH_API_KEY)")

        self.timeout = timeout

        if base_url:
            self.base_url = base_url.rstrip("/")
        else:
            self.base_url = SANDBOX_URL if environment == "sandbox" else PRODUCTION_URL

        self.session = requests.Session()
        self.session.headers.update({
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json",
            "User-Agent": "Respatch-Python-SDK/1.0.0",
        })

    def _request(self, method: str, path: str, **kwargs) -> Any:
        url = f"{self.base_url}{path}"

        try:
            response = self.session.request(method, url, timeout=self.timeout, **kwargs)

            if response.status_code >= 400:
                try:
                    error_data = response.json()
                    message = error_data.get("message", response.text)
                except Exception:
                    message = response.text
                raise RespatchAPIError(message, status_code=response.status_code, response=response)

            if response.status_code == 204:
                return None

            return unwrap_data(response.json())
        except requests.RequestException as e:
            raise RespatchAPIError(str(e)) from e

    def get(self, path: str, params: Optional[Dict[str, Any]] = None) -> Any:
        return self._request("GET", path, params=params)

    def post(self, path: str, json: Optional[Dict[str, Any]] = None) -> Any:
        return self._request("POST", path, json=json)

    def put(self, path: str, json: Optional[Dict[str, Any]] = None) -> Any:
        return self._request("PUT", path, json=json)

    def delete(self, path: str) -> Any:
        return self._request("DELETE", path)


TransmitClient = RespatchClient
