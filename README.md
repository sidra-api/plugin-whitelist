# Plugin Whitelist

## Description  
The Whitelist Plugin filters incoming requests based on their IP addresses in the Sidra Api. This plugin ensures that only IPs included in the whitelist are allowed to access the backend services.

---

## How It Works  

1. **Initializing Allowed IPs**  
   - The list of allowed IPs is retrieved from the `ALLOWED_IPS` environment variable.  
   - If the environment variable is not set, the plugin uses the following default list:  
     ```
     192.168.1.1,192.168.1.2
     ```

2. **IP Address Validation**  
   - The plugin checks the client's IP from the following headers:  
     - `X-Real-Ip`  
     - `X-Forwarded-For`  
     - `Remote-Addr`  
   - If the client's IP matches any of the whitelisted IPs, access is granted.

3. **Response**  
   - If the IP is allowed:  
     - **Status**: `200 OK`  
     - **Body**: "Allowed".  
   - If the IP is not allowed:  
     - **Status**: `403 Forbidden`  
     - **Body**: "IP not allowed".

---

## Installation

### Clone Repository  
```bash
git clone <repository-url>
cd plugin-whitelist
