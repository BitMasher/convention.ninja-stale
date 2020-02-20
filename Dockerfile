# build front end
FROM node:current-alpine as frontendbuilder
WORKDIR /workspace
COPY ./frontend/ frontend/
WORKDIR /workspace/frontend
RUN yarn build

# build backend
FROM golang:alpine as backendbuilder
WORKDIR /workspace
COPY ./backend/ backend/
WORKDIR /workspace/backend
RUN go build -o server --ldflags "-w -s" main.go

# execution environment
FROM alpine:latest
WORKDIR /server
COPY --from=backendbuilder /workspace/backend/server /server/server
COPY --from=frontendbuilder /workspace/frontend/build/ /server/static/
ENV PORT=3000
WORKDIR /server
CMD ["/server/server"]