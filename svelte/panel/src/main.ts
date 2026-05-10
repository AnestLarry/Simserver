import "./app.pcss";
import "./app.css";
import App from "./App.svelte";
import { mount } from "svelte";
import { init } from '$lib/i18n/i18n.svelte'
init()

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;
