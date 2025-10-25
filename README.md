# Conditional Headers Plugin for Traefik

Set request headers conditionally based on the incoming hostname. Supports multiple hosts sharing the same headers.

## Features

- Multiple hosts can share the same header configuration
- Supports exact host matching
- Supports wildcard subdomain matching (e.g., `*.demo.dev.io`)
- First matching rule wins

## Configuration

### Multiple hosts with same headers:
```yaml
http:
  middlewares:
    my-conditional-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - demo.dev.io
                - test.demo.dev.io
                - staging.demo.dev.io
              headers:
                X-App-Kernel-Name: demo.valksor.com
                X-Environment: production
            - hosts:
                - prod.example.com
              headers:
                X-App-Kernel-Name: prod.example.com
                X-Environment: production
```

### With wildcards:
```yaml
http:
  middlewares:
    my-conditional-headers:
      plugin:
        conditional-headers:
          rules:
            - hosts:
                - "*.internal.dev.io"
              headers:
                X-App-Kernel-Name: internal.valksor.com
            - hosts:
                - "*.public.dev.io"
              headers:
                X-App-Kernel-Name: public.valksor.com
```

## Docker Compose Labels
```yaml
labels:
  - "traefik.http.middlewares.my-headers.plugin.conditional-headers.rules[0].hosts[0]=demo.dev.io"
  - "traefik.http.middlewares.my-headers.plugin.conditional-headers.rules[0].hosts[1]=test.demo.dev.io"
  - "traefik.http.middlewares.my-headers.plugin.conditional-headers.rules[0].hosts[2]=staging.demo.dev.io"
  - "traefik.http.middlewares.my-headers.plugin.conditional-headers.rules[0].headers.X-App-Kernel-Name=demo.valksor.com"
```
