# traffic-clones-api

An apiserver for <https://shields.io/endpoint>.
Get the total number of clones and the total number of unique clones of your github repo.

## Description

Provides:

- traffic clones total count api:

  `/v1/github/traffic/clones/total?git_user=&git_repo=&type=count&label=clones-total`

- traffic clones total uniques api:

  `/v1/github/traffic/clones/total?git_user=&git_repo=&type=uniques&label=clones-uniques`

## Deployment Server

Login your cloud server(like aws EC2) first.

#### 1. Generate a github personal access token

Access this page <https://github.com/settings/tokens>,
and press button `Generate new token` to generate your access token.

For example, the generated access token is `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`.

#### 2. Deploy backend api server `traffic-clones-api`

```sh
$ git clone --depth 1 https://github.com/windvalley/traffic-clones-api

$ cd traffic-clones-api

$ go build

# it will listening on :9000 port
$ nohup ./traffic-clones-api -t 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 &
```

#### 3. Deploy https proxy server by [Caddy](https://github.com/caddyserver/caddy)

Assume that your domain name `api.sre.im` has been resolved to your cloud server IP.

```sh
$ wget https://github.com/caddyserver/caddy/releases/download/v2.4.6/caddy_2.4.6_linux_amd64.tar.gz

$ tar zxf caddy_2.4.6_linux_amd64.tar.gz

$ mv ./caddy /usr/local/bin/

$ mkdir api.sre.im

$ cd api.sre.im

$ cat > Caddyfile <<EOF
{
    http_port 80
    https_port 443
}

api.sre.im  {
    reverse_proxy 127.0.0.1:9000

    log {
        output file logs/access.log
        format single_field common_log
    }
}
EOF

$ nohup caddy run &
```

#### 4. Test the deployment

Request the url `https://api.sre.im/v1/github/traffic/clones/total?git_user=your-github-username&git_repo=your-github-repo&type=count-or-uniques&label=clones-xxx`

The response should be:

```json
{
  "schemaVersion": 1,
  "label": "clones-xxx",
  "message": "1728",
  "color": "orange"
}
```

## Create the badge

#### 1. Open <https://shields.io/endpoint> in browser explorer

#### 2. Generate the final badge url

Add the url `https://api.sre.im/v1/github/traffic/clones/total?git_user=your-github-username&git_repo=your-github-repo&type=count-or-uniques&label=clones-xxx` in `url` blank.

Then click `Copy Badge URL` to copy it to system clipboard.

Finally, add your badge link `![traffic clones total](the content of system clipboard)` to the `README.md` of your github repo.

The badge will be shown as:

[![Traffic Clones Total](https://img.shields.io/endpoint?url=https%3A%2F%2Fapi.sre.im%2Fv1%2Fgithub%2Ftraffic%2Fclones%2Ftotal%3Fgit_user%3Dwindvalley%26git_repo%3Dgossh%26type%3Dcount%26label%3Dclones-total)](https://github.com/windvalley/traffic-clones-api)

## License

This project is under the MIT License.
See the [LICENSE](LICENSE) file for the full license text.
