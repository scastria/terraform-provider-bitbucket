# Bitbucket Provider
The Bitbucket provider is used to interact with the Bitbucket API.  The provider
needs to be configured with the proper credentials before it can be used.

This provider does NOT cover 100% of the Bitbucket API.  If there is something missing
that you would like to be added, please submit an Issue in corresponding GitHub repo.
## Example Usage
```hcl
terraform {
  required_providers {
    bitbucket = {
      source  = "scastria/bitbucket"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Bitbucket Provider
provider "bitbucket" {
  client_id = "XXXX"
  client_secret = "YYYY"
  num_retries = 3
  retry_delay = 30
}
```
## Argument Reference
* `access_token` - **(Optional, String)** The access token obtained via authentication that can be used instead of `client_id` and `client_secret`. Token Authentication. Can be specified via env variable `BB_ACCESS_TOKEN`.
* `client_id` - **(Optional, String)** The client_id that will invoke all Bitbucket API commands. Client Credentials Authentication. Can be specified via env variable `BB_CLIENT_ID`.
* `client_secret` - **(Optional, String)** The client_secret for the client_id. Client Credentials Authentication. Can be specified via env variable `BB_CLIENT_SECRET`.
* `num_retries` - **(Optional, Integer)** Number of retries for each Bitbucket API call in case of 429-Too Many Requests or any 5XX status code. Can be specified via env variable `BB_NUM_RETRIES`. Default: 3.
* `retry_delay` - **(Optional, Integer)** How long to wait (in seconds) in between retries. Can be specified via env variable `BB_RETRY_DELAY`. Default: 30.
