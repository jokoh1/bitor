<script lang="ts">
    import { onMount } from 'svelte';
    import { pocketbase } from '$lib/stores/pocketbase';
    import { Card, Chart, Heading } from 'flowbite-svelte';
    import type { ApexOptions } from 'apexcharts';
  
    let chartOptions: ApexOptions = {
      chart: {
        type: 'bar',
        height: 350,
      },
      series: [
        {
          name: 'Vulnerabilities',
          data: [0, 0, 0, 0], // Default data
        },
      ],
      xaxis: {
        categories: ['Critical', 'High', 'Medium', 'Unknown'],
      },
      plotOptions: {
        bar: {
          horizontal: false,
          columnWidth: '50%',
          endingShape: 'rounded',
          distributed: true, // Add this line
        },
      },
      dataLabels: {
        enabled: true,
      },
      yaxis: {
        title: {
          text: 'Number of Vulnerabilities',
        },
      },
      fill: {
        opacity: 1,
      },
      tooltip: {
        y: {
          formatter: function (val) {
            return val.toString();
          },
        },
      },
      colors: ['#8B5CF6', '#DC2626', '#FACC15', '#6B7280'], // Aligned colors
    };
  
    let vulnerabilitiesData = [];
  
    async function fetchOpenVulnerabilities() {
      try {
        const result = await $pocketbase.collection('nuclei_results').getFullList(200, {
          filter: `
            (severity = "critical" || severity = "high" || severity = "medium" || severity = "unknown")
            && acknowledged = false
            && false_positive = false
          `,
        });
  
        // Process the data to count vulnerabilities by severity
        const severityCounts = {
          Critical: 0,
          High: 0,
          Medium: 0,
          Unknown: 0,
        };
  
        result.forEach((finding) => {
          const severity = finding.severity.toLowerCase();
          switch (severity) {
            case 'critical':
              severityCounts.Critical += 1;
              break;
            case 'high':
              severityCounts.High += 1;
              break;
            case 'medium':
              severityCounts.Medium += 1;
              break;
            case 'unknown':
              severityCounts.Unknown += 1;
              break;
            // Ignore 'low' and 'info' severities
          }
        });
  
        vulnerabilitiesData = [
          severityCounts.Critical,
          severityCounts.High,
          severityCounts.Medium,
          severityCounts.Unknown,
        ];
  
        // Update chart data
        chartOptions.series[0].data = vulnerabilitiesData;
      } catch (error) {
        console.error('Error fetching open vulnerabilities:', error);
      }
    }
  
    onMount(() => {
      fetchOpenVulnerabilities();
    });
  </script>
  
  <Card size="xl" class="w-full max-w-none 2xl:col-span-2">
    <div class="mb-4 flex items-center justify-between">
      <div class="flex-shrink-0">
        <Heading tag="h3" class="text-2xl">All Open Vulnerabilities</Heading>
        <p class="text-base font-light text-gray-500 dark:text-gray-400">
          Vulnerabilities not acknowledged or marked as false positive
        </p>
      </div>
    </div>
  
    {#if chartOptions.series[0].data.some((value) => value > 0)}
      <Chart options={chartOptions}></Chart>
    {:else}
      <p class="text-center text-gray-500 dark:text-gray-400">No data available to display.</p>
    {/if}
  </Card>