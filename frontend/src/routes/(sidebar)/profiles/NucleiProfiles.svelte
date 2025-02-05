<script lang="ts">
    import { onMount } from 'svelte';
    import {
      Button,
      Table,
      TableHead,
      TableHeadCell,
      TableBody,
      TableBodyRow,
      TableBodyCell,
      Heading,
      Modal,
      Input,
      Label,
      Alert
    } from 'flowbite-svelte';
    import Card from '../../utils/Card.svelte';
    import { pocketbase } from '$lib/stores/pocketbase';
    import AddProfileForm from './AddProfileForm.svelte';
    import ProfileBuilderModal from './ProfileBuilderModal.svelte';
  
    const defaultConfig = {
      'max-host-error': 500,
      'bulk-size': 2000,
      'concurrency': 32,
      'error-log': 'nuclei-errors.log',
      'stats': true,
      'no-color': true
    };
  
    let showProfileForm = false;
    let showProfileBuilder = false;
    let currentProfileData: Record<string, any> = {
      name: '',
      profile: '',
      id: '',
      raw_yaml: ''
    };
    let profiles: Array<any> = [];
  
    async function createDefaultProfile() {
      try {
        // Check if default profile already exists
        const existingDefault = await $pocketbase.collection('nuclei_profiles').getFirstListItem(`name="Default Profile"`).catch(() => null);
        
        if (!existingDefault) {
          await $pocketbase.collection('nuclei_profiles').create({
            name: 'Default Profile',
            profile: defaultConfig,
            raw_yaml: `# Default Nuclei Configuration
# These settings provide a balanced configuration for most scans

max-host-error: 500
bulk-size: 2000
concurrency: 32
error-log: nuclei-errors.log
stats: true
no-color: true`,
            is_default: true
          });
          
          await fetchProfiles();
        }
      } catch (error) {
        console.error('Error creating default profile:', error);
      }
    }
  
    async function fetchProfiles() {
      try {
        const result = await $pocketbase.collection('nuclei_profiles').getList();
        profiles = result.items;
      } catch (error) {
        console.error('Error fetching profiles:', error);
      }
    }
  
    onMount(async () => {
      await createDefaultProfile();
      await fetchProfiles();
    });
  
    function openAddProfileModal() {
      currentProfileData = { name: '', profile: '', id: '', raw_yaml: '' };
      showProfileForm = true;
    }
  
    function openEditModal(profile) {
      currentProfileData = { ...profile };
      showProfileForm = true;
    }
  
    async function deleteProfile(id: string) {
      if (confirm('Are you sure you want to delete this profile?')) {
        try {
          await $pocketbase.collection('nuclei_profiles').delete(id);
          fetchProfiles();
        } catch (error) {
          console.error('Error deleting profile:', error);
        }
      }
    }
  </script>
  
<main class="p-4">
  <Card class="w-full">
    <div class="space-y-8">
      <!-- Heading for Nuclei Profiles -->
      <Heading tag="h2" class="text-xl font-semibold mb-4">Nuclei Profiles</Heading>

      <!-- Buttons to Add Profile -->
      <div class="flex space-x-4 mb-4">
        <Button on:click={openAddProfileModal}>Add Profile (YAML)</Button>
        <Button on:click={() => (showProfileBuilder = true)}>Profile Builder</Button>
      </div>

      <!-- Profiles Table -->
      <div>
        <Table class="w-full border border-gray-200 dark:border-gray-700">
          <TableHead class="bg-gray-100 dark:bg-gray-700">
            {#each ['Name', 'Created', 'Actions'] as title}
              <TableHeadCell class="ps-4 font-normal">{title}</TableHeadCell>
            {/each}
          </TableHead>
          <TableBody>
            {#each profiles as profile}
              <TableBodyRow class="text-base hover:bg-gray-50 dark:hover:bg-gray-800">
                <TableBodyCell class="p-4">{profile.name}</TableBodyCell>
                <TableBodyCell class="p-4">{new Date(profile.created).toLocaleDateString()}</TableBodyCell>
                <TableBodyCell class="space-x-2">
                  <Button size="sm" on:click={() => openEditModal(profile)}>Edit</Button>
                  <Button size="sm" color="failure" on:click={() => deleteProfile(profile.id)}>Delete</Button>
                </TableBodyCell>
              </TableBodyRow>
            {/each}
          </TableBody>
        </Table>
      </div>

      <!-- Add/Edit Profile Modal -->
      {#if showProfileForm}
        <AddProfileForm
          bind:open={showProfileForm}
          bind:currentProfileData
          on:refreshProfiles={fetchProfiles}
        />
      {/if}

      <!-- Profile Builder Modal -->
      {#if showProfileBuilder}
        <ProfileBuilderModal bind:open={showProfileBuilder} />
      {/if}
    </div>
  </Card>
</main>