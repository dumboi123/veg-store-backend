#!/bin/bash

create_env_file() {
    # Clean old content before overwrite
    : > docker/.env

    echo JWT_ALG=$JWT_ALG >> docker/.env
    echo CERT_SECRET=$CERT_SECRET >> docker/.env

    echo POSTGRES_USER=$POSTGRES_USER >> docker/.env
    echo POSTGRES_PASSWORD=$POSTGRES_PASSWORD >> docker/.env

    echo GIN_MODE=$GIN_MODE >> docker/.env
    echo MODE=$MODE >> docker/.env

}

generate_keypair() {
    local keypair_dir="$SECRETS_DIR/keypair"

    # === 📂 Prepare ca folder ===
    if [[ -d "$keypair_dir" ]]; then
      rm -rf "$keypair_dir"/*
    else
      mkdir -p "$keypair_dir"
    fi

    local private_key="$keypair_dir/private.pem"
    local public_key="$keypair_dir/public.pem"

    openssl genpkey -algorithm RSA -out "$private_key" -pkeyopt rsa_keygen_bits:2048
    openssl rsa -pubout -in "$private_key" -out "$public_key"
    echo "Keypair are generated."
}
