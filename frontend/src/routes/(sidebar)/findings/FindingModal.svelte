<script lang="ts">
  import { Modal, Button, Toggle } from 'flowbite-svelte';
  import { Accordion, AccordionItem } from 'flowbite-svelte';
  import { pocketbase } from '$lib/stores/pocketbase';
  import { createEventDispatcher, onDestroy, onMount } from 'svelte';
  // Import Xterm components and types
  import { Xterm, XtermAddon } from '@battlefieldduck/xterm-svelte';
  import type { ITerminalOptions, ITerminalInitOnlyOptions, Terminal } from '@battlefieldduck/xterm-svelte';
  import BlockNoteEditor from '$lib/BlockNoteEditor.svelte';

  let showNotes = false;
  let saveMessage = '';
  let copyMessage = '';

  let notesContent: string = '';
  let previousFindingId = '';
  let editorKey = 0;
  const dispatch = createEventDispatcher();

  export let open = false;

  // Flag to track if user has toggled showNotes
  let hasUserToggled = false;

  interface Finding {
    id: string;
    info: {
      name: string;
      tags: string[];
      description?: string;
      reference?: string | string[];
    };
    host: string;
    ip: string;
    template_id: string;
    severity: 'critical' | 'high' | 'medium' | 'low' | 'info' | 'unknown';
    severity_order: number;
    client: {
      name: string;
    } | null;
    scan: {
      name: string;
      id: string;
    } | null;
    timestamp: string;
    last_seen: string;
    acknowledged: boolean;
    false_positive: boolean;
    request?: string;
    response?: string;
    notes?: string;
    updated?: string;
    matched_at?: string;
    type?: string;
    extra_info?: any;
    port?: number;
    url?: string;
    curl_command?: string;
    remediated: boolean;
  }

  export let finding: Finding | null = null;

  // Terminal options
  const terminalOptions: ITerminalOptions & ITerminalInitOnlyOptions = {
    fontFamily: 'Consolas, Monaco, "Courier New", monospace',
    fontSize: 14,
    theme: {
      background: '#1e1e1e',
      foreground: '#ffffff',
    },
    cursorBlink: true,
    disableStdin: true,
    convertEol: true,
  };

  // Terminal instances
  let requestTerminalInstance: Terminal | null = null;
  let responseTerminalInstance: Terminal | null = null;
  let curlTerminalInstance: Terminal | null = null;

  let requestFitAddon: any;
  let responseFitAddon: any;
  let curlFitAddon: any;
  // Reactive variable to hold parsed references
  let parsedReferences: string[] = [];

    // Reactive statement to detect finding changes
    $: if (finding && finding.id !== previousFindingId) {
    console.log('Finding changed from', previousFindingId, 'to', finding.id);
    previousFindingId = finding.id;
    fetchNotes().then(() => {
      editorKey += 1; // Increment key to force re-render
    });
  }

  // Update parsedReferences when finding changes
  $: if (finding) {
    if (typeof finding.info.reference === 'string') {
      try {
        const urlRegex = /https?:\/\/[^\s"]+/g;
        parsedReferences = finding.info.reference.match(urlRegex) || [];
      } catch (error) {
        console.error('Error parsing references from string:', error);
      }
    } else if (Array.isArray(finding.info.reference)) {
      parsedReferences = finding.info.reference;
    } else {
      parsedReferences = [];
    }

    // Initialize notes content if available
    if (finding.notes) {
      notesContent = finding.notes;
    } else {
      notesContent = '';
    }
  }

  async function fetchNotes() {
    if (finding) {
      try {
        console.log('Fetching notes for finding ID:', finding.id);
        const record = await $pocketbase.collection('nuclei_results').getOne(finding.id);
        console.log('Fetched record:', record);
        if (record.notes) {
          notesContent = record.notes;
          if (!hasUserToggled) {
            showNotes = true;
          }
        } else {
          notesContent = '';
          if (!hasUserToggled) {
            showNotes = false;
          }
        }
        console.log('Updated notesContent:', notesContent);
      } catch (error) {
        console.error('Error fetching notes:', error);
        notesContent = '';
        if (!hasUserToggled) {
          showNotes = false;
        }
      }
    }
  }

  // Function to initialize terminals with content
  async function onCurlLoad(event: CustomEvent<{ terminal: Terminal }>) {
    curlTerminalInstance = event.detail.terminal;

    // Load the FitAddon
    curlFitAddon = new (await XtermAddon.FitAddon()).FitAddon();
    curlTerminalInstance.loadAddon(curlFitAddon);
    curlFitAddon.fit();

    // Write request data to the terminal
    if (finding && finding.curl_command) {
      curlTerminalInstance.reset();
      curlTerminalInstance.write(finding.curl_command);
    } else {
      console.log('No curl command available');
    }
  }

  

  // Function to initialize terminals with content
  async function onRequestLoad(event: CustomEvent<{ terminal: Terminal }>) {
    requestTerminalInstance = event.detail.terminal;

    // Load the FitAddon
    requestFitAddon = new (await XtermAddon.FitAddon()).FitAddon();
    requestTerminalInstance.loadAddon(requestFitAddon);
    requestFitAddon.fit();

    // Write request data to the terminal
    if (finding && finding.request) {
      requestTerminalInstance.reset();
      requestTerminalInstance.write(finding.request);
    } else {
      console.log('No request data available');
    }
  }

  async function onResponseLoad(event: CustomEvent<{ terminal: Terminal }>) {
    responseTerminalInstance = event.detail.terminal;

    // Load the FitAddon
    responseFitAddon = new (await XtermAddon.FitAddon()).FitAddon();
    responseTerminalInstance.loadAddon(responseFitAddon);
    responseFitAddon.fit();

    // Write response data to the terminal
    if (finding && finding.response) {
      responseTerminalInstance.reset();
      responseTerminalInstance.write(finding.response);
    } else {
      console.log('No response data available');
    }
  }

  // Refit terminals when modal opens
  $: if (open) {
    if (requestFitAddon) {
      requestFitAddon.fit();
    }
    if (responseFitAddon) {
      responseFitAddon.fit();
    }
  }

  // Function to handle acknowledged toggle
  async function toggleAcknowledged() {
    if (finding) {
      try {
        const updatedFinding = await $pocketbase.collection('nuclei_results').update(finding.id, {
          acknowledged: finding.acknowledged,
        });
        console.log('Acknowledged status updated:', updatedFinding);
        saveMessage = 'Acknowledged status updated successfully';
      } catch (error) {
        console.error('Error updating acknowledged status:', error);
        saveMessage = 'Error updating acknowledged status';
      }
    }
  }

  // Function to handle false positive toggle
  async function toggleFalsePositive() {
    if (finding) {
      try {
        const updatedFinding = await $pocketbase.collection('nuclei_results').update(finding.id, {
          false_positive: finding.false_positive,
        });
        console.log('False positive status updated:', updatedFinding);
        saveMessage = 'False positive status updated successfully';
      } catch (error) {
        console.error('Error updating false positive status:', error);
        saveMessage = 'Error updating false positive status';
      }
    }
  }

  // Function to handle remediated toggle
  async function toggleRemediated() {
    if (finding) {
      try {
        const updatedFinding = await $pocketbase.collection('nuclei_results').update(finding.id, {
          remediated: finding.remediated,
        });
        console.log('Remediated status updated:', updatedFinding);
        saveMessage = 'Remediated status updated successfully';
      } catch (error) {
        console.error('Error updating remediated status:', error);
        finding.remediated = !finding.remediated; // Revert on error
        saveMessage = 'Error updating remediated status';
      }
    }
  }

  // Function to handle content updates from the editor
  function handleNotesChange(newContent: string) {
    notesContent = newContent;
    // Additional logic if needed
  }

  // Function to handle manual toggle by user
  function handleShowNotesToggle() {
    hasUserToggled = true;
  }

  // Update notesContent and reset toggle when finding changes
  $: if (finding && finding.id !== previousFindingId) {
    previousFindingId = finding.id;
    hasUserToggled = false;
    fetchNotes();
  }

  // Function to save notes
  async function saveNotes() {
    if (finding) {
      try {
        // Save the notesContent back to the backend
        await $pocketbase.collection('nuclei_results').update(finding.id, {
          notes: notesContent,
        });
        // Update the local finding object
        finding.notes = notesContent;
        saveMessage = 'Notes saved successfully';
        setTimeout(() => {
          saveMessage = '';
        }, 2000);
      } catch (error) {
        console.error('Error saving notes:', error);
        saveMessage = 'Error saving notes';
        setTimeout(() => {
          saveMessage = '';
        }, 2000);
      }
    }
  }

  function closeModal() {
    open = false;
    // Reset messages upon closing
    saveMessage = '';
    copyMessage = '';
    showNotes = false;
  }

  // Get severity color
  function getSeverityColor(severity: string): string {
      switch (severity.toLowerCase()) {
        case 'critical':
          return 'bg-red-600 text-white';
        case 'high':
          return 'bg-orange-500 text-white';
        case 'medium':
          return 'bg-yellow-500 text-white';
        case 'low':
          return 'bg-green-500 text-white';
        case 'info':
          return 'bg-blue-500 text-white';
        default:
          return 'bg-gray-500 text-white';
    }
  }

  // Function to copy text to clipboard
  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text).then(
      () => {
        copyMessage = 'Copied to clipboard!';
        setTimeout(() => {
          copyMessage = '';
        }, 2000);
      },
      (err) => {
        console.error('Could not copy text: ', err);
        copyMessage = 'Failed to copy text.';
        setTimeout(() => {
          copyMessage = '';
        }, 2000);
      }
    );
  }

  // Add after the reactive statement
  $: {
    console.log('Finding ID:', finding?.id);
    console.log('showNotes:', showNotes);
    console.log('hasUserToggled:', hasUserToggled);
  }

  function handleChange(event: CustomEvent<string>) {
    console.log('Editor content changed:', event.detail);
    notesContent = event.detail;
  }

  // Add this helper function for formatting dates
  function formatDate(dateString: string | undefined): string {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleString();
  }

  let staleThresholdDays = 30;

  onMount(async () => {
    try {
      // Get stale threshold from system settings
      const systemSettings = await $pocketbase.collection('system_settings').getFirstListItem('');
      if (systemSettings) {
        staleThresholdDays = systemSettings.stale_threshold_days || 30;
      }
    } catch (error) {
      console.error('Error loading system settings:', error);
    }
  });

  // Update the helper function to use timestamp if last_seen is not available
  function getLastSeenStatus(lastSeen: string | undefined, timestamp: string): { isStale: boolean; daysAgo: number } {
    // Use lastSeen if available, otherwise use timestamp
    const dateToCheck = new Date(lastSeen || timestamp);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - dateToCheck.getTime());
    const daysAgo = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return {
      isStale: daysAgo > staleThresholdDays,
      daysAgo
    };
  }

</script>

<Modal bind:open={open} size="xl" placement="center" dismissible={true} keyboard={false}>
  <div class="relative bg-white rounded-lg shadow dark:bg-gray-700">
    <!-- Modal header -->
    <div class="flex items-start justify-between p-4 border-b rounded-t dark:border-gray-600">
      <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
        Finding Details
      </h3>
    </div>
    <!-- Modal body -->
    <div class="p-6 space-y-6">
      {#if copyMessage}
        <div class="copy-message-global">
          {copyMessage}
        </div>
      {/if}

      {#if finding}
        <!-- Display finding details -->
        <div>
          <p><strong>Name:</strong> {finding.info.name}</p>
          <p>
            <strong>Severity:</strong>
            <span
              class={`inline-block px-2 py-1 rounded ${getSeverityColor(finding.severity)}`}
            >
              {finding.severity}
            </span>
          </p>
          <!-- Display tags -->
          {#if finding.info.tags && finding.info.tags.length}
            <p><strong>Tags:</strong> {finding.info.tags.join(', ')}</p>
          {/if}
          <!-- Display description -->
          {#if finding.info.description}
            <p><strong>Description:</strong></p>
            <p>{@html finding.info.description}</p>
          {/if}
          <!-- Display references -->
          {#if parsedReferences.length}
            <p><strong>References:</strong></p>
            <ul>
              {#each parsedReferences as reference}
                <li><a href={reference} target="_blank" rel="noopener noreferrer">{reference}</a></li>
              {/each}
            </ul>
          {/if}

          <!-- Additional fields -->
          {#if finding.extra_info}
            <p><strong>Extra Info:</strong> {JSON.stringify(finding.extra_info)}</p>
          {/if}

          {#if finding.port}
            <p><strong>Port:</strong> {finding.port}</p>
          {/if}

          {#if finding.url}
            <p><strong>URL:</strong> {finding.url}</p>
          {/if}

          {#if finding.curl_command}
            <Accordion flush={true} class="mt-4">
              <AccordionItem>
                <span slot="header">Curl Command</span>
                <div class="terminal-container relative">
                  <!-- Copy Button -->
                  <Button
                    class="absolute top-2 right-2 z-10"
                    size="xs"
                    on:click={() => copyToClipboard(finding.curl_command)}
                  >
                    Copy
                  </Button>

                  <!-- Copy Message Notification -->
                  {#if copyMessage}
                    <div class="copy-message">
                      {copyMessage}
                    </div>
                  {/if}

                  <Xterm
                    options={terminalOptions}
                    on:load={onCurlLoad}
                  />
                </div>
              </AccordionItem>
            </Accordion>
          {/if}

          <!-- Accordion for Request and Response -->
          {#if finding.request || finding.response}
            <Accordion flush={true} class="mt-4">
              {#if finding.request}
                <AccordionItem>
                  <span slot="header">Request</span>
                  <div class="terminal-container relative">
                    <!-- Copy Button -->
                    <Button
                      class="absolute top-2 right-2 z-10"
                      size="xs"
                      on:click={() => copyToClipboard(finding.request)}
                    >
                      Copy
                    </Button>
                    <Xterm
                      options={terminalOptions}
                      on:load={onRequestLoad}
                    />
                  </div>
                </AccordionItem>
              {/if}
              {#if finding.response}
                <AccordionItem>
                  <span slot="header">Response</span>
                  <div class="terminal-container relative">
                    <!-- Copy Button -->
                    <Button
                      class="absolute top-2 right-2 z-10"
                      size="xs"
                      on:click={() => copyToClipboard(finding.response)}
                    >
                      Copy
                    </Button>
                    <Xterm
                      options={terminalOptions}
                      on:load={onResponseLoad}
                    />
                  </div>
                </AccordionItem>
              {/if}
            </Accordion>
          {/if}

          <!-- Acknowledged Toggle -->
          <div class="flex items-center space-x-2 mt-4">
            <Toggle bind:checked={finding.acknowledged} on:change={toggleAcknowledged} />
            <span class="text-gray-700 dark:text-gray-300">Acknowledged</span>
          </div>

          <!-- False Positive Toggle -->
          <div class="flex items-center space-x-2">
            <Toggle bind:checked={finding.false_positive} on:change={toggleFalsePositive} />
            <span class="text-gray-700 dark:text-gray-300">False Positive</span>
            {#if finding.false_positive}
              <span class="text-red-500 text-xl ml-2">😔</span>
            {:else}
              <span class="text-green-500 text-xl ml-2">😊</span>
            {/if}
          </div>

          <!-- Remediated Toggle -->
          <div class="flex items-center space-x-2">
            <Toggle bind:checked={finding.remediated} on:change={toggleRemediated} />
            <span class="text-gray-700 dark:text-gray-300">Remediated</span>
            {#if finding.remediated}
              <span class="text-green-500 text-xl ml-2">✅</span>
            {:else}
              <span class="text-yellow-500 text-xl ml-2">⚠️</span>
            {/if}
          </div>

          <!-- Notes Section -->
          <div class="mt-4">
            <div class="flex items-center space-x-2">
              <Toggle bind:checked={showNotes} on:change={handleShowNotesToggle} />
              <span class="text-gray-700 dark:text-gray-300">Add Notes</span>
            </div>

            {#if showNotes}
              <div class="mt-2">
                <BlockNoteEditor 
                  content={notesContent} 
                  on:change={handleChange} 
                />
                <Button class="mt-2" on:click={saveNotes}>Save Notes</Button>
                {#if saveMessage}
                  <div class="text-green-500 mt-2">{saveMessage}</div>
                {/if}
              </div>
            {/if}
          </div>

          <div class="space-y-2">
            <div class="flex flex-col space-y-1">
              <span class="text-sm font-medium text-gray-500">First Seen</span>
              <span class="text-sm">{finding.timestamp ? formatDate(finding.timestamp) : 'N/A'}</span>
              
              {#if finding.last_seen && finding.last_seen !== finding.timestamp}
                <span class="text-sm font-medium text-gray-500 mt-2">Last Seen</span>
                <span class="text-sm">{formatDate(finding.last_seen)}</span>
                {#if getLastSeenStatus(finding.last_seen, finding.timestamp).isStale}
                  <div class="mt-1 text-amber-600 dark:text-amber-500 flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    <span>Not seen in {getLastSeenStatus(finding.last_seen, finding.timestamp).daysAgo} days</span>
                  </div>
                {/if}
              {:else if getLastSeenStatus(undefined, finding.timestamp).isStale}
                <div class="mt-1 text-amber-600 dark:text-amber-500 flex items-center gap-1">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <span>Only seen once, {getLastSeenStatus(undefined, finding.timestamp).daysAgo} days ago</span>
                </div>
              {/if}
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</Modal>

<style>
  .terminal-container {
    position: relative;
    height: 200px;
    overflow: auto;
    background-color: #1e1e1e;
    border-radius: 5px;
    padding: 10px;
  }
  /* Additional styles if needed */

  .copy-message {
    position: absolute;
    top: 40px; /* Adjust as needed */
    right: 10px;
    z-index: 10;
    background-color: #d1fae5; /* Light green background */
    color: #065f46; /* Dark green text */
    padding: 5px 10px;
    border-radius: 5px;
    font-size: 0.875rem; /* text-sm */
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .copy-message-global {
    background-color: #d1fae5; /* Light green background */
    color: #065f46; /* Dark green text */
    padding: 10px;
    border-radius: 5px;
    font-size: 1rem;
    margin-bottom: 10px;
    text-align: center;
  }
</style>