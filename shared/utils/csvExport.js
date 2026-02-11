/**
 * Converts an array of objects to CSV format
 * @param {Array<Object>} data - Array of data objects
 * @param {string} filename - Filename for the download
 */
export function exportToCSV(data, filename = 'export.csv') {
  if (!data || data.length === 0) {
    console.warn('No data to export');
    return;
  }

  // Get all unique keys from all objects (handles varying fields)
  const allKeys = new Set();
  data.forEach(row => {
    Object.keys(row).forEach(key => allKeys.add(key));
  });
  
  const headers = Array.from(allKeys);

  // Convert headers to CSV row
  const csvHeaders = headers.map(header => escapeCSVValue(header)).join(',');

  // Convert data rows to CSV
  const csvRows = data.map(row => {
    return headers.map(header => {
      const value = row[header];
      
      // Handle NULL/undefined as empty string
      if (value === null || value === undefined) {
        return '';
      }
      
      // Handle dates
      if (value instanceof Date) {
        return escapeCSVValue(value.toISOString());
      }
      
      // Handle objects (JSON stringify them)
      if (typeof value === 'object') {
        return escapeCSVValue(JSON.stringify(value));
      }
      
      // Convert to string and escape
      return escapeCSVValue(String(value));
    }).join(',');
  });

  // Combine headers and rows
  const csv = [csvHeaders, ...csvRows].join('\n');

  // Create blob and download
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  
  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  URL.revokeObjectURL(url);
}

/**
 * Escapes a value for CSV format
 * @param {string} value - Value to escape
 * @returns {string} Escaped value
 */
function escapeCSVValue(value) {
  if (value === null || value === undefined) {
    return '';
  }
  
  const stringValue = String(value);
  
  // If value contains comma, quote, or newline, wrap in quotes and escape quotes
  if (stringValue.includes(',') || stringValue.includes('"') || stringValue.includes('\n')) {
    return `"${stringValue.replace(/"/g, '""')}"`;
  }
  
  return stringValue;
}

/**
 * Generate a filename with timestamp
 * @param {string} prefix - Filename prefix
 * @param {string} extension - File extension (default: csv)
 * @returns {string} Filename with timestamp
 */
export function generateFilename(prefix = 'export', extension = 'csv') {
  const now = new Date();
  const timestamp = now.toISOString().replace(/[:.]/g, '-').slice(0, -5);
  return `${prefix}_${timestamp}.${extension}`;
}
