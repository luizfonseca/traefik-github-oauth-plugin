# Traefik GitHub OAuth Plugin

> This is a fork of [MuXiu1997](https://github.com/MuXiu1997/traefik-github-oauth-plugin) repository. This fork is mostly fixing some of the security concerns I wanted to address. This will be kept synced with the main repo.

This is a Traefik middleware plugin that allows users to authenticate using GitHub OAuth.

The plugin is intended to be used as a replacement for the BasicAuth middleware,

providing a more secure way for users to access protected routes.

![process](https://user-images.githubusercontent.com/49554020/216764214-4097f8da-33d2-49ef-9f12-0194d671bd92.svg)

## Quick Start (Docker)

1. Create a GitHub OAuth App

   - See: https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app
   - Set the Authorization callback URL to `http://<traefik-github-oauth-server-host>/oauth/redirect`

2. Run the Traefik GitHub OAuth server

   ```sh
   docker run -d --name traefik-github-oauth-server \
     --network <traefik-proxy-network> \
     -e 'GITHUB_OAUTH_CLIENT_ID=<client-id>' \
     -e 'GITHUB_OAUTH_CLIENT_SECRET=<client-secret>' \
     -e 'API_BASE_URL=http://<traefik-github-oauth-server-host>' \
     -l 'traefik.http.services.traefik-github-oauth-server.loadbalancer.server.port=80' \
     -l 'traefik.http.routers.traefik-github-oauth-server.rule=Host(`<traefik-github-oauth-server-host>`)' \
     luizfonseca/traefik-github-oauth-server
   ```

3. Install the Traefik GitHub OAuth plugin

    Add this snippet in the Traefik Static configuration

   ```yaml
   experimental:
     plugins:
       github-oauth:
         moduleName: "github.com/luizfonseca/traefik-github-oauth-plugin"
         version: <version>
   ```

4. Run your App

   ```sh
   docker run -d --whoami test \
     --network <traefik-proxy-network> \
     --label 'traefik.http.middlewares.whoami-github-oauth.plugin.github-oauth.apiBaseUrl=http://traefik-github-oauth-server' \
     --label 'traefik.http.middlewares.whoami-github-oauth.plugin.github-oauth.whitelist.logins[0]=luizfonseca' \
     --label 'traefik.http.middlewares.whoami-github-oauth.plugin.github-oauth.whitelist.teams[0]=827726' \
     --label 'traefik.http.middlewares.whoami-github-oauth.plugin.github-oauth.whitelist.twoFactorAuthRequired=true' \
     --label 'traefik.http.routers.whoami.rule=Host(`whoami.example.com`)' \
     --label 'traefik.http.routers.whoami.middlewares=whoami-github-oauth' \
    traefik/whoami
   ```

## Configuration

### Server configuration

| Environment Variable         | Description                                                                   | Default | Required |
|------------------------------|-------------------------------------------------------------------------------|---------|----------|
| `GITHUB_OAUTH_CLIENT_ID`     | The GitHub OAuth App client id                                                |         | Yes      |
| `GITHUB_OAUTH_CLIENT_SECRET` | The GitHub OAuth App client secret                                            |         | Yes      |
| `GITHUB_OAUTH_SCOPES`        | Additional scopes to be added to the Oauth workflow.                          |         | No       |
| `API_BASE_URL`               | The base URL of the Traefik GitHub OAuth server                               |         | Yes      |
| `API_SECRET_KEY`             | The api secret key. You can ignore this if you are using the internal network |         | No       |
| `SERVER_ADDRESS`             | The server address                                                            | `:80`   | No       |
| `DEBUG_MODE`                 | Enable debug mode and set log level to debug                                  | `false` | No       |
| `LOG_LEVEL`                  | The log level, Available values: debug, info, warn, error                     | `info`  | No       |

### Middleware Configuration

```yaml
# The base URL of the Traefik GitHub OAuth server
apiBaseUrl: http://<traefik-github-oauth-server-host>
# The api secret key. You can ignore this if you are using the internal network
apiSecretKey: optional_secret_key_if_not_on_the_internal_network
# The path to redirect to after the user has authenticated, defaults to /_auth
# Note: This path is not GitHub OAuth App's Authorization callback URL
authPath: /_auth
# optional jwt secret key, if not set, the plugin will generate a random key
jwtSecretKey: optional_secret_key
# The log level, defaults to info
# Available values: debug, info, warn, error
logLevel: info

# whitelist
whitelist:
  # When set to `true`, the middleware will check if the given user has 2FA
  # configured, otherwise they will be denied access
  # Default is `false`
  twoFactorAuthRequired: 'true'

  # The list of GitHub user ids that are whitelisted to access the resources
  ids:
    - 996

  # The list of GitHub user logins that are whitelisted to access the resources
  logins:
    - luizfonseca

  # The list of Github Teams that are whitelisted to access the resources
  teams:
    - 988772
```

### OAuth Configuration

For the OAuth configuration, you need to create a GitHub OAuth App.
You can follow the steps in the [GitHub documentation](https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app) to create it and obtain the `GITHUB_OAUTH_CLIENT_ID` and `GITHUB_OAUTH_CLIENT_SECRET` values.

#### OAuth Scopes
- For `ids` and `logins` you don't need extra scopes.
- For `teams` you might need to request the `read:org` scope from the user. See the [documentation](https://docs.github.com/en/rest/teams/teams?apiVersion=2022-11-28#list-teams-for-the-authenticated-user).
    - You can do so by updating the `GITHUB_OAUTH_SCOPES` environment variable with the desired additional scopes, e.g. `GITHUB_OAUTH_SCOPES="read:org"` via the **Server Configuration**.


## License

[MIT](./LICENSE)
