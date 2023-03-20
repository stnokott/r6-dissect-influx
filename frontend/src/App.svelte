<script lang="ts">
  import "carbon-components-svelte/css/g90.css";

  import { Button } from "carbon-components-svelte";
  import SettingsIcon from "carbon-icons-svelte/lib/Settings.svelte";

  import SettingsDialog from "./settings/Settings.svelte";
  import MatchList from "./matches/MatchList.svelte";
  import HeaderView from "./Header.svelte";
  import FooterView from "./footer/Footer.svelte";

  import {
    IsConfigComplete,
    Connect,
    Disconnect,
  } from "./../wailsjs/go/main/App";
  import type { ConnectionDetails } from "./db";
  import { onMount } from "svelte";
  import type { MatchListAPI } from "./matches/matchlist";

  let settingsDialogOpen: boolean;

  let isConfigComplete = false;
  let promConnectionDetails: Promise<ConnectionDetails>;

  let matchListAPI: MatchListAPI;

  $: {
    isConfigComplete;
    if (matchListAPI) {
      matchListAPI.onConfigChanged();
    }
  }

  $: if (isConfigComplete) {
    Disconnect().then(() => {
      promConnectionDetails = Connect();
    });
  }

  function openSettings(): void {
    settingsDialogOpen = true;
  }

  function onConfigChanged(_e: CustomEvent<void>) {
    // can assume that config is complete since only then can the config dialog be closed/confirmed
    isConfigComplete = true;
  }

  function onConnected(e: CustomEvent<ConnectionDetails>) {
    promConnectionDetails = new Promise((r) => r(e.detail));
  }

  onMount(async () => {
    isConfigComplete = await IsConfigComplete();
  });
</script>

<SettingsDialog
  bind:open={settingsDialogOpen}
  on:connected={onConnected}
  on:changed={onConfigChanged}
/>
<HeaderView />
<div id="root">
  <div id="content">
    {#if isConfigComplete}
      <MatchList bind:matchListAPI />
    {:else}
      <div class="container-centered">
        <Button kind="primary" icon={SettingsIcon} on:click={openSettings}
          >Setup</Button
        >
      </div>
    {/if}
  </div>
  <div id="footer">
    <FooterView on:openSettings={openSettings} {promConnectionDetails} />
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

    overflow: hidden auto;
    position: relative;
  }

  .container-centered {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);

    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: row;
  }

  #footer {
    flex-grow: 0;
  }
</style>
