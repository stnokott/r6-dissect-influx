<script lang="ts">
  // TODO: change to SCSS
  import "carbon-components-svelte/css/g90.css";

  import SettingsDialog from "./settings/Settings.svelte";
  import ContentView from "./Content.svelte";
  import HeaderView from "./Header.svelte";
  import FooterView from "./Footer.svelte";

  let settingsDialogOpen: boolean;

  function openSettings(): void {
    settingsDialogOpen = true;
  }
</script>

<SettingsDialog bind:open={settingsDialogOpen} />
<HeaderView />
<div id="root">
  <div id="content">
    <ContentView />
  </div>
  <div id="footer">
    <FooterView {openSettings} />
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
