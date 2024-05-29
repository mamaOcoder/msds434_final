FROM golang:latest
WORKDIR /app
COPY go.mod go.sum  ./
RUN go mod tidy
COPY . .
#COPY gcp_creds.json ./gcp_creds.json
#ENV GOOGLE_APPLICATION_CREDENTIALS=/app/gcp_creds.json
EXPOSE 8080
EXPOSE 9100
RUN go build -o main .
CMD [ "/app/main" ]