# Traefik Plugin for Treblle

## What you get with the plugin
This plugin allows you to integrate your traefik proxy with Treblle seamlessly so you can enjoy:
- Top-Notch Observability: Gain comprehensive insights into your API's operations.
- Auto-Generated API Documentation: Automatically keep your documentation up-to-date with your API changes.
- Alerting and Notifications: Stay informed with alerts on crucial API events and changes.

## Prerequisites
- Traefik: This plugin has been tested on v2 and v3
- [Treblle](https://app.treblle.com/) Account

## Installation
- Obtain your API key and Project ID from [Treblle](https://app.treblle.com/)
- Visit the Traefik [Plugin Catalog](https://plugins.traefik.io/plugins)
- Search for `Treblle` and select Install.
- Follow the displayed instructions to configure the plugin in your static and dynamic configurations.
- Restart your Traefik instance
- Enjoy!

## Configuration
| Parameter                  | Type   | Description                             |
|----------------------------|--------|-----------------------------------------|
| `ApiKey`                   | string | Your API key obtained from Treblle.     |
| `ProjectId`                | string | Your Project ID from Treblle.           |
| `AdditionalFieldsToMask`   | array  | Additional sensitive fields to mask.    |
| `RoutesToBlock`            | array  | Paths to routes that should be hidden.  |
| `RoutesRegex`              | string | Regex to match and hide specific routes.|

## Example Usage
```yaml
# Static configuration
experimental:
  plugins:
    treblle:
      moduleName: "github.com/Treblle/TreblleTraefikPluginGo"
      version: "{version-from-github-releases}"

# Dynamic configuration
http:
  routers:
    my-router:
      entryPoints:
        - http
      middlewares:
        - my-plugin
      service: service-whoami
      rule: Host(`localhost`)

  services:
    service-whoami:
      loadBalancer:
        servers:
          - url: "http://localhost:8081"
        passHostHeader: true

  middlewares:
    my-plugin:
      plugin:
        treblle:
          ApiKey: "your-api-key"
          ProjectId: "your-project-id"
          AdditionalFieldsToMask:
            - "accessToken"
            - "refreshToken"
          RoutesToBlock:
            - "/api/user/login"
            - "/ping"
          RoutesRegex: "^/api/projects"
```