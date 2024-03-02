FROM node:18 AS build-stage

WORKDIR /app

RUN npm install -g pnpm

COPY package.json ./
COPY pnpm-lock.yaml ./
RUN pnpm install
COPY . ./
RUN npx vite build --mode=staging
RUN ls

FROM nginx:1.24-alpine AS deploy-stage
COPY --from=build-stage /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
