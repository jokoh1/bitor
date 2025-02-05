<script>
  import { Modal } from 'flowbite-svelte';

  export let finding = null;
  export let open = false;

  function closeModal() {
    open = false;
  }
</script>

<Modal bind:open={open} size="xl" placement="center">
  <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
    <h3 class="text-lg leading-6 font-medium text-gray-900">
      Finding Details
    </h3>
    <div class="mt-2 space-y-2">
      <p><strong>ID:</strong> {finding.id}</p>
      <p><strong>Collection ID:</strong> {finding.collectionId}</p>
      <p><strong>Collection Name:</strong> {finding.collectionName}</p>
      <p><strong>Created:</strong> {finding.created}</p>
      <p><strong>Updated:</strong> {finding.updated}</p>
      <p><strong>Template ID:</strong> {finding.template_id}</p>
      <p><strong>Scan:</strong> {finding.scan}</p>
      <p><strong>Name:</strong> {finding.name}</p>
      <p><strong>Description:</strong> {finding.description}</p>
      <p><strong>Severity:</strong> {finding.severity}</p>
      <p><strong>Matched At:</strong> {finding.matched_at}</p>
      <p><strong>Type:</strong> {finding.type}</p>
      <p><strong>Host:</strong> {finding.host}</p>
      <p><strong>IP:</strong> {finding.ip}</p>
      <p><strong>Timestamp:</strong> {finding.timestamp}</p>
      <p><strong>Reference:</strong> {finding.reference}</p>
      <p><strong>Extra Info:</strong> {JSON.stringify(finding.extra_info)}</p>
      <p><strong>Info:</strong> {JSON.stringify(finding.info)}</p>
      <p><strong>Port:</strong> {finding.port}</p>
      <p><strong>URL:</strong> {finding.url}</p>
      <p><strong>Request:</strong> {finding.request}</p>
      <p><strong>Response:</strong> {finding.response}</p>
      <p><strong>Matcher Status:</strong> {finding.matcher_status ? 'True' : 'False'}</p>
      <p><strong>Tags:</strong> {JSON.stringify(finding.tags)}</p>
    </div>
  </div>
  <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
    <button type="button" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" on:click={closeModal}>
      Close
    </button>
  </div>
</Modal>
