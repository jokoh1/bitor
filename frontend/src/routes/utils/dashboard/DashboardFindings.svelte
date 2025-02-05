<script lang="ts">
  import { onMount } from 'svelte';
  import { pocketbase } from '$lib/stores/pocketbase';
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Card,
    Heading,
    Button
  } from 'flowbite-svelte';
  import FindingModal from './FindingModal.svelte';

  interface Finding {
    info: {
      name: string;
      tags: string[];
    };
    host: string;
    ip: string;
    template_id: string;
    severity: string;
    client: {
      name: string;
      favicon?: string;
    };
    timestamp: string;
  }

  let findings: Finding[] = [];
  let selectedFinding = null;
  let showModal = false;
  let currentPage = 1;
  const itemsPerPage = 10;
  let sortField = '';
  let sortDirection = 'asc';

  // Fetch data from PocketBase
  async function fetchFindings(page = 1, sortField = '', sortDirection = 'asc') {
    try {
      const result = await $pocketbase.collection('nuclei_results').getList(page, itemsPerPage, {
        sort: `${sortDirection === 'asc' ? '' : '-'}${sortField}`
      });
      const findingsWithClients = await Promise.all(
        result.items.map(async (item) => {
          const client = await $pocketbase.collection('clients').getOne(item.client);
          return {
            ...item,
            client: client
          };
        })
      );

      if (page === 1) {
        findings = findingsWithClients;
      } else {
        findings = findings.concat(findingsWithClients);
      }
    } catch (error) {
      console.error('Error fetching findings:', error);
    }
  }

  function sortData(field: string) {
    if (sortField === field) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortField = field;
      sortDirection = 'asc';
    }

    // Fetch sorted data from the server
    fetchFindings(currentPage, sortField, sortDirection);
  }

  function openModal(finding: Finding) {
    selectedFinding = finding;
    showModal = true;
  }

  function changePage(direction: string) {
    if (direction === 'next' && currentPage * itemsPerPage < findings.length) {
      currentPage++;
      fetchFindings(currentPage);
    } else if (direction === 'prev' && currentPage > 1) {
      currentPage--;
    }
  }

  function getSeverityColor(severity: string) {
    switch (severity.toLowerCase()) {
      case 'high':
        return 'bg-red-500 text-white';
      case 'medium':
        return 'bg-yellow-500 text-white';
      case 'low':
        return 'bg-green-500 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  }

  onMount(() => {
    fetchFindings();
  });
</script>

<Card size="xl" class="shadow-sm max-w-none">
  <Heading tag="h3" class="-ml-0.25 mb-2 text-xl font-semibold dark:text-white">
    Nuclei Findings
  </Heading>
  <Table
    hoverable={true}
    noborder
    striped
    class="mt-6 min-w-full divide-y divide-gray-200 dark:divide-gray-600"
  >
  <TableHead class="bg-gray-50 dark:bg-gray-700">
    <TableHeadCell on:click={() => sortData('info.name')}>
      Name {#if sortField === 'info.name'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell on:click={() => sortData('host')}>
      Host {#if sortField === 'host'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell on:click={() => sortData('ip')}>
      IP Address {#if sortField === 'ip'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell on:click={() => sortData('template_id')}>
      Template ID {#if sortField === 'template_id'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell on:click={() => sortData('severity')}>
      Severity {#if sortField === 'severity'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell on:click={() => sortData('info.tags')}>
      Tags {#if sortField === 'info.tags'}{sortDirection === 'asc' ? '▲' : '▼'}{/if}
    </TableHeadCell>
    <TableHeadCell>
        Client
    </TableHeadCell>
    <TableHeadCell>
        Timestamp
    </TableHeadCell>
  </TableHead>
    <TableBody>
      {#each findings.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage) as finding}
        <TableBodyRow on:click={() => openModal(finding)}>
          <TableBodyCell class="px-4 font-normal">{finding.info.name}</TableBodyCell>
          <TableBodyCell class="px-4 font-normal text-gray-500 dark:text-gray-400">
            {finding.host}
          </TableBodyCell>
          <TableBodyCell class="px-4">{finding.ip}</TableBodyCell>
          <TableBodyCell class="px-4 font-normal text-gray-500 dark:text-gray-400">
            {finding.template_id}
          </TableBodyCell>
          <TableBodyCell class="px-4 font-normal">
            <span class={`inline-block px-2 py-1 rounded ${getSeverityColor(finding.severity)}`}>
              {finding.severity}
            </span>
          </TableBodyCell>
          <TableBodyCell class="px-4 font-normal">
            {#each finding.info.tags as tag}
              <span class="inline-block bg-gray-200 text-gray-800 text-xs px-2 py-1 rounded-full mr-1">{tag}</span>
            {/each}
          </TableBodyCell>
          <TableBodyCell class="px-4 font-normal">
            {#if finding.client}
              <div class="flex items-center gap-2">
                {#if finding.client.favicon}
                  <img src={$pocketbase.getFileUrl(finding.client, finding.client.favicon)} alt="{finding.client.name} Favicon" class="h-4 w-4" />
                {/if}
                {finding.client.name}
              </div>
            {:else}
              N/A
            {/if}
          </TableBodyCell>
          <TableBodyCell class="px-4 font-normal">
            {finding.timestamp || 'N/A'}
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  <div class="pagination-controls">
    <Button on:click={() => changePage('prev')} disabled={currentPage === 1}>Previous</Button>
    <Button on:click={() => changePage('next')}>Next</Button>
  </div>
</Card>

<FindingModal bind:open={showModal} finding={selectedFinding} />
