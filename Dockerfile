FROM oryd/oathkeeper:v0.40.9
ADD config.yaml /config.yaml
ADD rules.yaml /rules.yaml
ADD jwks.json /jwks.json
