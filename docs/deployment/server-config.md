# TeamSeeker Deployment Configuration

## Configuration Files Location

### 1. Nginx Configuration
**Location**: `/etc/nginx/sites-available/teamseeker`

```nginx
# First Server Block - DDNS Configuration
server {
    listen 38000;
    server_name pang1.ddns.net;
    
    # Frontend proxy
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # API proxy
    location /api/ {
        proxy_pass http://localhost:3000/api/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

# Second Server Block - Ngrok SSL Configuration
server {
    listen 443 ssl;
    server_name *.ngrok-free.app;
    
    # SSL Configuration
    ssl_certificate /etc/ssl/certs/nginx-selfsigned.crt;
    ssl_certificate_key /etc/ssl/private/nginx-selfsigned.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # Frontend proxy
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # API proxy
    location /api/ {
        proxy_pass http://localhost:3000/api/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. Systemd Service Configuration
**Location**: `/etc/systemd/system/teamseeker.service`

```ini
[Unit]
Description=TeamSeeker Service
After=network.target

[Service]
Type=simple
User=pang1
WorkingDirectory=/home/pang1/Dev/TeamSeekerBackend
ExecStart=/usr/local/bin/TeamSeeker
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

### 3. Ngrok Configuration

#### Ngrok Service
**Location**: `/etc/systemd/system/ngrok.service`
```ini
[Unit]
Description=ngrok
After=network.target

[Service]
Type=simple
User=pang1
Environment=NGROK_AUTHTOKEN=your_auth_token_here
ExecStart=/usr/local/bin/ngrok start --config=/home/pang1/.ngrok2/ngrok.yml backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

#### Ngrok YAML Config
**Location**: `/home/pang1/.ngrok2/ngrok.yml`
```yaml
version: "2"
tunnels:
  backend:
    addr: http://localhost:38000
    proto: http
    hostname: reasonably-proven-quetzal.ngrok-free.app
```

## Port Configuration
- Nginx DDNS: Port 38000
- Nginx SSL: Port 443
- Application: Port 3000
- Ngrok: Forwarding from port 38000

## SSL Certificates
- Location: `/etc/ssl/certs/nginx-selfsigned.crt`
- Private Key: `/etc/ssl/private/nginx-selfsigned.key`
- Protocols: TLSv1.2, TLSv1.3

## Quick Commands

### Service Management
```bash
# Start TeamSeeker
sudo systemctl start teamseeker

# Stop TeamSeeker
sudo systemctl stop teamseeker

# Restart TeamSeeker
sudo systemctl restart teamseeker

# Check status
sudo systemctl status teamseeker

# View logs
sudo journalctl -u teamseeker -f
```

### Nginx Management
```bash
# Test configuration
sudo nginx -t

# Reload configuration
sudo systemctl reload nginx

# Restart Nginx
sudo systemctl restart nginx

# Check Nginx status
sudo systemctl status nginx
```

### Ngrok Management
```bash
# Start/stop/restart ngrok service
sudo systemctl start ngrok
sudo systemctl stop ngrok
sudo systemctl restart ngrok

# Check ngrok service status
sudo systemctl status ngrok

# View ngrok logs
sudo journalctl -u ngrok -f

# Check ngrok tunnels
curl localhost:4040/api/tunnels
```

## Directory Structure
```
/home/pang1/
├── Dev/
│   └── TeamSeekerBackend/  # Application directory
└── .ngrok2/
    └── ngrok.yml           # Ngrok configuration

/etc/
├── nginx/
│   └── sites-available/
│       └── teamseeker      # Nginx configuration
├── systemd/system/
│   └── teamseeker.service  # Systemd service
└── ssl/
    ├── certs/
    │   └── nginx-selfsigned.crt
    └── private/
        └── nginx-selfsigned.key
```

## Security Notes
1. The ngrok authtoken is stored as an environment variable in the systemd service file
2. When updating the ngrok authtoken:
   - Edit the ngrok service file: `sudo nano /etc/systemd/system/ngrok.service`
   - Update the `NGROK_AUTHTOKEN` environment variable
   - Reload systemd: `sudo systemctl daemon-reload`
   - Restart ngrok: `sudo systemctl restart ngrok`

## Important Notes
1. The application runs under the user 'pang1'
2. The backend binary is located at `/usr/local/bin/TeamSeeker`
3. Two server configurations:
   - DDNS access via port 38000
   - Ngrok SSL access via port 443
4. All traffic is proxied to localhost:3000
