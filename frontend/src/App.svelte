<script lang="ts">
  // TODO: change to SCSS
  import "carbon-components-svelte/css/g90.css";

  import SettingsDialog from "./settings/Settings.svelte";
  import ContentView from "./Content.svelte";
  import HeaderView from "./Header.svelte";
  import FooterView from "./Footer.svelte";

  import type { db } from "./settings/types";

  let settingsDialogOpen: boolean;
  let connectionDetails: db.ConnectionDetails;

  function openSettings(): void {
    settingsDialogOpen = true;
  }

  function onConnected(details: db.ConnectionDetails): void {
    connectionDetails = details;
  }
</script>

<SettingsDialog bind:open={settingsDialogOpen} {onConnected} />
<HeaderView />
<div id="root">
  <div id="content">
    <ContentView />
  </div>
  <div id="footer">
    <FooterView {openSettings} {connectionDetails} />
  </div>
</div>

<style>
  #root {
    height: calc(100vh - 4rem);

    display: flex;
    flex-direction: column;
  }

  #content {
    flex-grow: 1;
    padding: 1em;

    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  #footer {
    flex-grow: 0;
  }
</style>
