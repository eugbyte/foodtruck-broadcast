# About

UI in `svelte` to track geographical data of foodtrucks.

The rendering strategy is using [SSG](https://www.educative.io/answers/ssr-vs-csr-vs-isr-vs-ssg).

# dev

```
npm i -g pnpm   // this project uses `pnpm`
pnpm i
pnpm run dev
```

## PWA mode

The service worker is bundled for production, but not during development

```
pnpm run build
pnpm run preview
```

Another to thing to note is that since the rendering strategy is SSG which deploys a new HTML file for each route, each route `+layout.sevelte` need to reference to the `manifest.webmanifest` file.

# deploy

```
pnpm run build
```

# test

Unit test

```
pnpm run test:unit
```

E2E test

```
pnpm run test:e2e
```
