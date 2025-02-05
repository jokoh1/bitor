<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { BlockNoteEditor } from '@blocknote/core';
  import '@blocknote/core/style.css';

  let editor: BlockNoteEditor;
  let editorElement: HTMLElement;
  let menuElement: HTMLDivElement;
  let isMenuVisible = false;
  let menuPosition = { top: 0, left: 0 };

  type MenuItem = {
    title: string;
    onItemClick: () => void;
    group?: string;
    icon?: string;
  };

  function createMenuItems(): MenuItem[] {
    return [
      // Basic blocks
      {
        title: 'Text',
        group: 'Basic blocks',
        icon: '📝',
        onItemClick: () => updateBlockType('paragraph')
      },
      {
        title: 'Heading 1',
        group: 'Basic blocks',
        icon: 'H1',
        onItemClick: () => updateBlockType('heading', { level: 1 })
      },
      {
        title: 'Heading 2',
        group: 'Basic blocks',
        icon: 'H2',
        onItemClick: () => updateBlockType('heading', { level: 2 })
      },
      {
        title: 'Heading 3',
        group: 'Basic blocks',
        icon: 'H3',
        onItemClick: () => updateBlockType('heading', { level: 3 })
      },
      // Lists
      {
        title: 'Bullet List',
        group: 'Lists',
        icon: '•',
        onItemClick: () => updateBlockType('bulletListItem')
      },
      {
        title: 'Numbered List',
        group: 'Lists',
        icon: '1.',
        onItemClick: () => updateBlockType('numberedListItem')
      },
      {
        title: 'Check List',
        group: 'Lists',
        icon: '☐',
        onItemClick: () => updateBlockType('checkListItem')
      },
      // Advanced
      {
        title: 'Code Block',
        group: 'Advanced',
        icon: '💻',
        onItemClick: () => updateBlockType('codeBlock')
      }
    ];
  }

  function updateBlockType(type: string, props: Record<string, any> = {}) {
    editor.updateBlock(editor.getTextCursorPosition().block, {
      type: type as any,
      props
    });
    isMenuVisible = false;
  }

  const menuItems = createMenuItems();

  onMount(() => {
    menuElement = document.createElement('div');
    menuElement.className = 'bn-suggestion-menu';
    document.body.appendChild(menuElement);

    editor = BlockNoteEditor.create();
    editor.mount(editorElement);

    editorElement.addEventListener('keydown', (event) => {
      if (event.key === '/') {
        const selection = editor._tiptapEditor.view.state.selection;
        const pos = editor._tiptapEditor.view.coordsAtPos(selection.from);
        menuPosition = { top: pos.top + 20, left: pos.left };
        isMenuVisible = true;
        event.preventDefault();
      } else if (event.key === 'Escape') {
        isMenuVisible = false;
      }
    });
  });

  onDestroy(() => {
    if (editor) {
      editor.mount(null);
    }
    if (menuElement) {
      menuElement.remove();
    }
  });

  $: if (isMenuVisible && menuElement) {
    menuElement.style.display = 'block';
    menuElement.style.top = `${menuPosition.top}px`;
    menuElement.style.left = `${menuPosition.left}px`;
    menuElement.innerHTML = '';

    // Group items
    const groups = menuItems.reduce((acc, item) => {
      const group = item.group || 'Other';
      if (!acc[group]) acc[group] = [];
      acc[group].push(item);
      return acc;
    }, {} as Record<string, MenuItem[]>);

    // Create menu HTML
    Object.entries(groups).forEach(([groupName, items]) => {
      const groupDiv = document.createElement('div');
      groupDiv.className = 'menu-group';
      
      const groupTitle = document.createElement('div');
      groupTitle.className = 'menu-group-title';
      groupTitle.textContent = groupName;
      groupDiv.appendChild(groupTitle);

      items.forEach(item => {
        const button = document.createElement('button');
        button.className = 'menu-item';
        button.innerHTML = `
          <span class="menu-item-icon">${item.icon || ''}</span>
          <span class="menu-item-title">${item.title}</span>
        `;
        button.onclick = item.onItemClick;
        groupDiv.appendChild(button);
      });

      menuElement.appendChild(groupDiv);
    });
  } else if (menuElement) {
    menuElement.style.display = 'none';
  }
</script>

<div class="editor-container">
  <div bind:this={editorElement} class="editor-instance"></div>
</div>

<style>
  .editor-container {
    position: relative;
    width: 100%;
    height: calc(100vh - 70px);
    padding: 1rem;
  }

  .editor-instance {
    border: 1px solid #ccc;
    border-radius: 4px;
    min-height: 200px;
  }

  :global(.bn-suggestion-menu) {
    position: fixed !important;
    z-index: 99999 !important;
    background: white !important;
    border: 1px solid #ccc !important;
    border-radius: 4px !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
    padding: 4px !important;
    min-width: 200px !important;
  }

  :global(.menu-group) {
    padding: 4px 0;
    border-bottom: 1px solid #eee;
  }

  :global(.menu-group:last-child) {
    border-bottom: none;
  }

  :global(.menu-group-title) {
    padding: 4px 12px;
    font-size: 0.8em;
    color: #666;
    font-weight: 500;
  }

  :global(.menu-item) {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 8px 12px;
    border: none;
    background: none;
    text-align: left;
    cursor: pointer;
    gap: 8px;
  }

  :global(.menu-item:hover) {
    background: #f5f5f5;
  }

  :global(.menu-item-icon) {
    width: 20px;
    text-align: center;
    font-size: 0.9em;
  }

  :global(.menu-item-title) {
    flex: 1;
  }
</style>