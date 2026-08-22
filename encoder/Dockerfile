
FROM golang:1.25.0-alpine3.11

# ============================================================
# Variáveis do Bento4
# ============================================================

ENV BENTO4_VERSION="1-6-0-641" \
    BENTO4_NAME="Bento4-SDK-1-6-0-641.x86_64-unknown-linux" \
    BENTO4_PATH="/opt/bento4" \
    BENTO4_BIN="/opt/bento4/bin"

ENV PATH="${PATH}:${BENTO4_BIN}"

# ============================================================
# Dependências do sistema
# ============================================================

RUN apk add --no-cache \
    ffmpeg \
    bash \
    curl \
    wget \
    unzip \
    make \
    ca-certificates

# ============================================================
# Instalação do Bento4
#
# O domínio antigo www.bento4.com/downloads não está mais
# servindo corretamente esse arquivo.
#
# O arquivo 1.6.0-641 está disponível em bok.net.
# ============================================================

WORKDIR /tmp/bento4

RUN wget -q \
    "https://www.bok.net/Bento4/binaries/${BENTO4_NAME}.zip" \
    -O bento4.zip && \
    unzip -q bento4.zip && \
    mkdir -p "${BENTO4_PATH}" && \
    cp -R "${BENTO4_NAME}"/* "${BENTO4_PATH}/" && \
    rm -rf \
        "${BENTO4_NAME}" \
        bento4.zip

# ============================================================
# Validação do Bento4
# ============================================================

RUN test -x "${BENTO4_BIN}/mp4fragment" && \
    test -x "${BENTO4_BIN}/mp4dash" && \
    test -x "${BENTO4_BIN}/mp4encrypt" && \
    echo "Bento4 instalado com sucesso:" && \
    ls -la "${BENTO4_BIN}"

# ============================================================
# Diretório da aplicação
# ============================================================

WORKDIR /go/src

# ============================================================
# Container de desenvolvimento
# ============================================================

ENTRYPOINT ["tail", "-f", "/dev/null"]

