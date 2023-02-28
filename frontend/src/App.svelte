<script lang="ts">
  import "carbon-components-svelte/css/g90.css";

  import { Button, Loading } from "carbon-components-svelte";
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
  import type { db } from "./index";
  import { onMount } from "svelte";
  import type { MatchListAPI } from "./matches/matchlist";

  let settingsDialogOpen: boolean;

  let promIsConfigComplete: Promise<boolean>;
  let promConnectionDetails: Promise<db.ConnectionDetails>;

  let matchListAPI: MatchListAPI;

  $: {
    if (promIsConfigComplete) {
      promIsConfigComplete.then((complete) => {
        if (matchListAPI) {
          matchListAPI.onConfigChanged();
        }
        if (complete) {
          Disconnect().then(() => {
            promConnectionDetails = Connect();
          });
        }
      });
    }
  }

  function openSettings(): void {
    settingsDialogOpen = true;
  }

  function onConfigChanged(_e: CustomEvent<void>) {
    // can assume that config is complete since only then can the config dialog be closed/confirmed
    promIsConfigComplete = new Promise((r) => r(true));
  }

  function onConnected(e: CustomEvent<db.ConnectionDetails>) {
    promConnectionDetails = new Promise((r) => r(e.detail));
  }

  onMount(() => {
    promIsConfigComplete = IsConfigComplete();
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
    {#await promIsConfigComplete}
      <div class="container-centered">
        <Loading withOverlay={false} small />
        <span style="margin-left: .333rem;">Initializing...</span>
      </div>
    {:then isComplete}
      {#if isComplete}
        <MatchList bind:matchListAPI />
      {:else}
        <div class="container-centered">
          <Button kind="primary" icon={SettingsIcon} on:click={openSettings}
            >Setup</Button
          >
        </div>
      {/if}
    {/await}
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
