# build front end
FROM node:current-alpine as frontendbuilder
WORKDIR /workspace/frontend
COPY ./frontend/package.json package.json
COPY ./frontend/package-lock.json package-lock.json
RUN npm install
COPY ./frontend/* ./
COPY ./frontend/public ./public
COPY ./frontend/src ./src
RUN npm run build

# build backend
FROM golang:alpine as backendbuilder
WORKDIR /workspace
COPY ./backend backend
WORKDIR /workspace/backend
RUN ls
RUN go build -o server --ldflags "-w -s"

# execution environment
FROM alpine:latest
WORKDIR /server
COPY --from=backendbuilder /workspace/backend/server /server/server
COPY --from=frontendbuilder /workspace/frontend/build/ /server/static/
ENV PORT=3000
WORKDIR /server
CMD ["/server/server"]