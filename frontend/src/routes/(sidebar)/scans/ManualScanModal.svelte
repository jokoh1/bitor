<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import { Modal, Button, Alert } from 'flowbite-svelte';
  import { InfoCircleSolid } from 'flowbite-svelte-icons';
  import { pocketbase } from '$lib/stores/pocketbase';

  export let open: boolean;
  const dispatch = createEventDispatcher();

  let file: File | null = null;
  let selectedClientId = '';
  let clients = [];
  let userToken = $pocketbase.authStore.token;
  let uploadProgress = 0;
  let showAlert = false;

  async function fetchClients() {
    try {
      const result = await $pocketbase.collection('clients').getFullList();
      clients = result;
    } catch (error) {
      console.error('Error fetching clients:', error);
    }
  }

  onMount(() => {
    fetchClients();
  });

  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      file = input.files[0];
    }
  }

  async function handleImport() {
    if (!file || !selectedClientId) {
      alert('Please select a file and a client.');
      return;
    }

    try {
      const scanId = generateScanId();
      const client = clients.find(c => c.id === selectedClientId);
      const clientName = client ? client.name : 'Unknown Client';

      const newScan = await $pocketbase.collection('nuclei_scans').create({
        name: scanId,
        client: selectedClientId,
        status: 'Manual',
      });

      const newScanId = newScan.id;

      const formData = new FormData();
      formData.append('file', file);
      formData.append('client_id', selectedClientId);
      formData.append('scan_id', newScanId);

      const xhr = new XMLHttpRequest();
      xhr.open('POST', 'http://localhost:8090/api/scan/import-scan-results', true);
      xhr.setRequestHeader('Authorization', `Bearer ${userToken}`);

      xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
          uploadProgress = (event.loaded / event.total) * 100;
        }
      };

      xhr.onload = () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          console.log('File uploaded successfully.');
          dispatch('import');
        } else {
          console.error('Failed to upload file');
        }
      };

      xhr.onerror = () => {
        console.error('Error uploading file');
      };

      xhr.send(formData);

      // Immediately close the modal and show the alert
      open = false;
      showAlert = true;
    } catch (error) {
      console.error('Error during file upload:', error);
    }
  }

  function generateScanId() {
    return 'scan_manual_json_' + Date.now();
  }
</script>

<Modal bind:open={open} on:close={() => open = false} closeButton>
  <div slot="header">
    <h3 class="text-lg font-medium text-gray-900 dark:text-white">
      Manually Add Scan Results
    </h3>
  </div>
  <div class="flex flex-col gap-4 p-4">
    <div>
      <label for="jsonFile" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">
        Upload JSON File
      </label>
      <input
        type="file"
        id="jsonFile"
        accept=".json"
        on:change={handleFileChange}
        class="block w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 cursor-pointer"
      />
    </div>
    <div>
      <label for="clientSelect" class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">
        Select Client
      </label>
      <select
        id="clientSelect"
        bind:value={selectedClientId}
        class="block w-full p-2 text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300"
      >
        <option value="" disabled>Select a client</option>
        {#each clients as client}
          <option value={client.id}>{client.name}</option>
        {/each}
      </select>
    </div>
    <div class="w-full bg-gray-200 rounded-full h-2.5 dark:bg-gray-700">
      <div class="bg-blue-600 h-2.5 rounded-full" style="width: {uploadProgress}%"></div>
    </div>
  </div>
  <div slot="footer">
    <Button on:click={handleImport}>Import</Button>
    <Button color="gray" on:click={() => open = false}>Cancel</Button>
  </div>
</Modal>

{#if showAlert}
  <Alert color="green" on:close={() => showAlert = false}>
    <InfoCircleSolid slot="icon" class="w-5 h-5" />
    <span class="font-medium">Success alert!</span> File uploaded successfully. The server is processing the results.
  </Alert>
{/if}
