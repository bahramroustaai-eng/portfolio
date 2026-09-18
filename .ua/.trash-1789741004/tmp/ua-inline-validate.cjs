#!/usr/bin/env node
const fs = require('fs');
const graphPath = process.argv[2];
const outputPath = process.argv[3];
try {
  const graph = JSON.parse(fs.readFileSync(graphPath, 'utf8'));
  const issues = [], warnings = [];
  if (!Array.isArray(graph.nodes)) issues.push('graph.nodes is missing or not an array');
  if (!Array.isArray(graph.edges)) issues.push('graph.edges is missing or not an array');
  const nodes = Array.isArray(graph.nodes) ? graph.nodes : [];
  const edges = Array.isArray(graph.edges) ? graph.edges : [];
  const ids = new Set();
  const fileTypes = new Set(['file','config','document','service','pipeline','table','schema','resource','endpoint']);
  for (const [i, n] of nodes.entries()) {
    if (!n.id) issues.push(`Node[${i}] missing id`);
    if (ids.has(n.id)) issues.push(`Duplicate node ID '${n.id}'`);
    ids.add(n.id);
    for (const field of ['type','name','summary','tags']) if (!n[field] || (field === 'tags' && !n.tags.length)) issues.push(`Node[${i}] '${n.id}' missing ${field}`);
  }
  for (const [i, e] of edges.entries()) {
    if (!ids.has(e.source)) issues.push(`Edge[${i}] source '${e.source}' not found`);
    if (!ids.has(e.target)) issues.push(`Edge[${i}] target '${e.target}' not found`);
  }
  const assigned = new Set();
  if (!Array.isArray(graph.layers)) issues.push('graph.layers is missing or not an array');
  for (const layer of (graph.layers || [])) {
    for (const id of (layer.nodeIds || [])) {
      if (!ids.has(id)) issues.push(`Layer '${layer.id}' refs missing node '${id}'`);
      if (assigned.has(id)) issues.push(`Node '${id}' appears in multiple layers`);
      assigned.add(id);
    }
    for (const field of ['id','name','description','nodeIds']) if (!(field in layer)) issues.push(`Layer missing ${field}`);
  }
  for (const n of nodes.filter(n => fileTypes.has(n.type))) if (!assigned.has(n.id)) issues.push(`File node '${n.id}' not in any layer`);
  if (!Array.isArray(graph.tour)) issues.push('graph.tour is missing or not an array');
  for (const [i, step] of (graph.tour || []).entries()) {
    for (const field of ['order','title','description','nodeIds']) if (!(field in step)) issues.push(`Tour step ${i} missing ${field}`);
    for (const id of (step.nodeIds || [])) if (!ids.has(id)) issues.push(`Tour step ${i} refs missing node '${id}'`);
  }
  for (const n of nodes) if (!edges.some(e => e.source === n.id || e.target === n.id)) warnings.push(`Node '${n.id}' has no edges`);
  const stats = { totalNodes:nodes.length, totalEdges:edges.length, totalLayers:(graph.layers||[]).length, tourSteps:(graph.tour||[]).length, nodeTypes:nodes.reduce((a,n)=>(a[n.type]=(a[n.type]||0)+1,a),{}), edgeTypes:edges.reduce((a,e)=>(a[e.type]=(a[e.type]||0)+1,a),{}) };
  fs.writeFileSync(outputPath, JSON.stringify({issues,warnings,stats}, null, 2));
} catch (err) { process.stderr.write(err.message + '\n'); process.exit(1); }
