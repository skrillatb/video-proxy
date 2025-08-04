# video-proxy

Petit serveur Go permettant de proxifier des fichiers MP4 ou des vidéos hébergées sur Uqload depuis un serveur distant.

Utiliser un VPS disposant d’un excellent quota de bande passante (minimum 20 To) et d’au moins 40 Go de stockage.
Installer Go et Nginx si nécessaire.
Si les ressources sont insuffisantes, prendre un autre VPS du même type et répartir le trafic entre les deux.

Proxifier le tout via Nginx en exposant uniquement les ports 80 et 443.
Utiliser Cloudflare (mode proxy activé) pour avoir leur CDN et des avantages comme cache, protection DDoS, etc...

#### Endpoints disponibles

#### Uqload
```bash
https://proxy.domain.to/uqload?url=https://uqload.cx/embed-tpeaeo14n2xx.html
```

#### MP4
```bash
https://proxy.domain.to/mp4?url=https://monsite.com/video/singe.mp4
```

## Déploiement

1. Build Linux :
```
GOOS=linux GOARCH=amd64 go build -o vproxy
```

2. Upload vers le VPS
```terminal
scp ./vproxy .env videoproxy:/home/ubuntu/
```

3. Connexion au VPS
```terminal
ssh videoproxy
```

4. Rendre le binaire exécutable
```terminal
cd /home/ubuntu
chmod +x vproxy
```

5. Lancer le serveur en arrière-plan
```terminal
nohup ./vproxy > vproxy.log 2>&1 &
```
