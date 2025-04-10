#!/bin/bash

# Script para generar certificados EC y convertirlos a Base64
# Creamos un directorio para los certificados
mkdir -p certs
cd certs

# Generamos la clave privada EC usando la curva P-256 (secp256r1)
openssl ecparam -name prime256v1 -genkey -noout -out ec_private.pem

# Extraemos la clave pública de la clave privada
openssl ec -in ec_private.pem -pubout -out ec_public.pem

echo "Certificados EC generados correctamente en el directorio 'certs'"

# Convertimos los certificados a Base64 para usarlos en el código
PRIVATE_KEY_BASE64=$(cat ec_private.pem | base64)
PUBLIC_KEY_BASE64=$(cat ec_public.pem | base64)

echo ""
echo "JWT_EC_PRIVATE_KEY_BASE64=\"$PRIVATE_KEY_BASE64\""
echo ""
echo "JWT_EC_PUBLIC_KEY_BASE64=\"$PUBLIC_KEY_BASE64\""
echo ""