<script lang="ts">
  import "@material-design-icons/font/filled.css";
  import "@material-design-icons/font/outlined.css";

  import ConfigDialog from "./Dialog.svelte";
  import ContentView from "./Content.svelte";
  import HeaderView from "./Header.svelte";
  import FooterView from "./Footer.svelte";

  let configDialogOpen: boolean;

  function openSettings(): void {
    configDialogOpen = true;
  }
</script>

<ConfigDialog bind:open={configDialogOpen} />
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
    padding-top: 48px;
    height: calc(100vh - 48px);

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
