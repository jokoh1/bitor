<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';
    import { Input, Label, Modal, Select, Button } from 'flowbite-svelte';
    import { uniqueNamesGenerator, adjectives, colors, animals } from 'unique-names-generator';
    import { pocketbase } from '$lib/stores/pocketbase';
    import { getFavicon } from '$lib/utils/favicon';
  
    export let open = false;
    export let client: {
      name: string;
      logo?: string;
      hidden_name: string;
      id: string;
      group?: ClientGroup | null;
      homepage?: string;
      favicon?: string;
    } | null = null;

    interface ClientGroup {
      id: string;
      name: string;
    }

    export let mode: 'add' | 'edit' = 'add';
    export let size: string = 'lg';
  
    let clientName = '';
    let hiddenName = '';
    let selectedGroupId = '';
    let homepage = '';
    let clientGroups: ClientGroup[] = [];
  
    const dispatch = createEventDispatcher();
  
    let isInitialized = false;
  
    onMount(async () => {
      try {
        const result = await $pocketbase.collection('client_groups').getList();
        clientGroups = result.items.map(group => ({
          id: group.id,
          name: group.name
        }));
      } catch (error) {
        console.error('Error fetching client groups:', error);
      }
    });
  
    $: if (open && !isInitialized) {
      if (client) {
        console.log('client object:', client);
        clientName = client.name || '';
        hiddenName = client.hidden_name || '';
        selectedGroupId = client.group?.id || '';
        homepage = client.homepage || '';
        console.log('selectedGroupId after setting:', selectedGroupId);
      } else {
        clientName = '';
        hiddenName = '';
        selectedGroupId = '';
        homepage = '';
        console.log('Adding new client - form fields reset');
      }
      isInitialized = true;
    }
  
    $: if (!open) {
      isInitialized = false;
    }
  
    function generateCodeName() {
      hiddenName = uniqueNamesGenerator({
        dictionaries: [adjectives, colors, animals],
        separator: '-',
        length: 3,
      });
    }
  
    async function updateFavicon(url: string) {
      try {
        console.log('Fetching favicon for homepage:', url);
        const faviconUrl = await getFavicon(url);
        console.log('Got favicon URL:', faviconUrl);
        
        if (faviconUrl) {
          console.log('Fetching favicon from URL:', faviconUrl);
          
          // Add headers to prevent CORS issues
          const response = await fetch(faviconUrl, {
            headers: {
              'Accept': 'image/png,image/*',
            },
            mode: 'cors'
          });
          
          if (!response.ok) {
            throw new Error(`Failed to fetch favicon: ${response.status} ${response.statusText}`);
          }
          
          // Log response headers
          console.log('Response headers:', {
            type: response.headers.get('content-type'),
            size: response.headers.get('content-length')
          });
          
          // Get the response as a blob directly
          const blob = await response.blob();
          if (!blob) {
            throw new Error('Failed to get blob from response');
          }
          
          console.log('Got favicon blob:', {
            size: blob.size,
            type: blob.type
          });

          // Create a File object with the blob
          const file = new File([blob], 'favicon.png', { 
            type: 'image/png'
          });
          
          // Verify the file was created successfully
          if (!(file instanceof File)) {
            throw new Error('Failed to create File object');
          }
          
          console.log('Created favicon file:', {
            name: file.name,
            type: file.type,
            size: file.size
          });
          
          // Test reading the file to verify it's valid
          const reader = new FileReader();
          reader.onload = () => {
            console.log('Successfully read file contents');
          };
          reader.onerror = () => {
            console.error('Failed to read file contents');
          };
          reader.readAsArrayBuffer(file);
          
          return file;
        }
      } catch (error) {
        console.error('Error in updateFavicon:', error);
        return null;
      }
    }
  
    async function handleSave() {
      try {
        console.log('Starting save process...');
        const formData = new FormData();
        
        // Add basic fields
        formData.append('name', clientName);
        formData.append('hidden_name', hiddenName);
        if (selectedGroupId) {
          formData.append('group', selectedGroupId);
        }
        if (homepage) {
          formData.append('homepage', homepage);
        }

        // Handle favicon
        if (homepage) {
          console.log('Getting favicon for homepage:', homepage);
          const faviconFile = await updateFavicon(homepage);
          
          if (!faviconFile) {
            console.error('Failed to get favicon file');
          } else {
            console.log('Adding favicon to FormData:', {
              name: faviconFile.name,
              type: faviconFile.type,
              size: faviconFile.size
            });
            
            // Add the file to FormData
            formData.append('favicon', faviconFile, 'favicon.png');

            // Verify the file was added to FormData
            const formDataFile = formData.get('favicon') as File;
            if (!formDataFile) {
              console.error('Failed to add file to FormData');
            } else {
              console.log('Verified file in FormData:', {
                name: formDataFile.name,
                type: formDataFile.type,
                size: formDataFile.size
              });
            }
          }
        }

        // Log all FormData entries
        console.log('Final FormData contents:');
        for (const [key, value] of formData.entries()) {
          if (value instanceof File) {
            console.log(`${key}: File(${value.name}, ${value.type}, ${value.size} bytes)`);
          } else {
            console.log(`${key}: ${value}`);
          }
        }

        // Save to PocketBase
        let savedClient;
        if (mode === 'edit' && client) {
          console.log('Updating existing client...');
          savedClient = await $pocketbase.collection('clients').update(client.id, formData);
        } else {
          console.log('Creating new client...');
          savedClient = await $pocketbase.collection('clients').create(formData);
        }
        
        console.log('Raw saved client response:', savedClient);
        
        // Verify the favicon was saved
        if (savedClient.favicon) {
          const faviconUrl = $pocketbase.getFileUrl(savedClient, 'favicon');
          console.log('Saved favicon URL:', faviconUrl);
        } else {
          console.warn('No favicon in saved client. Raw favicon value:', savedClient.favicon);
        }
        
        dispatch('save', {
          name: clientName,
          hidden_name: hiddenName,
          group: selectedGroupId,
          homepage: homepage
        });
        
        open = false;
        clientName = '';
        hiddenName = '';
        selectedGroupId = '';
        homepage = '';
      } catch (error) {
        console.error('Save failed:', error);
        if (error && typeof error === 'object' && 'response' in error) {
          try {
            const errorData = await error.response.json();
            console.error('Error details:', errorData);
          } catch (e) {
            console.error('Could not parse error response');
          }
        }
      }
    }
  </script>

  <Modal bind:open size="lg" title={mode === 'edit' ? "Edit Client" : "Add New Client"}>
    <form on:submit|preventDefault={handleSave} class="space-y-6">
      <Label class="space-y-2">
        <span class="text-gray-700 dark:text-gray-300">Client Name</span>
        <Input bind:value={clientName} placeholder="Enter client name" required />
      </Label>
      <Label class="space-y-2">
        <span class="text-gray-700 dark:text-gray-300">Homepage</span>
        <div class="flex items-center space-x-2">
          <Input 
            bind:value={homepage} 
            placeholder="Enter company homepage (e.g., https://example.com)" 
            type="url"
          />
        </div>
      </Label>
      <Label class="space-y-2">
        <span class="text-gray-700 dark:text-gray-300">Hidden Name</span>
        <div class="flex items-center space-x-2">
          <Input bind:value={hiddenName} placeholder="Enter or generate hidden name" required />
          <Button on:click={generateCodeName} type="button">Generate</Button>
        </div>
      </Label>
      <Label class="space-y-2">
        <span class="text-gray-700 dark:text-gray-300">Client Group</span>
        <Select
          placeholder="Select a group"
          bind:value={selectedGroupId}
          on:change={() => console.log('Group changed to:', selectedGroupId)}
        >
          {#each clientGroups as group}
            <option value={group.id}>{group.name}</option>
          {/each}
        </Select>
      </Label>
      <Button type="submit" class="w-full text-white py-2">
        {mode === 'edit' ? "Update" : "Save"}
      </Button>
    </form>
  </Modal>