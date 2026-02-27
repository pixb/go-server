import { Timestamp } from '@bufbuild/protobuf';

/**
 * Convert protobuf Timestamp to JavaScript Date
 */
export function convertTimestampToDate(timestamp: Timestamp | undefined): Date | null {
  if (!timestamp) {
    return null;
  }
  return timestamp.toDate();
}

/**
 * Format protobuf Timestamp to display string
 * Format: YYYY-MM-DD HH:mm:ss
 */
export function formatTimestamp(timestamp: Timestamp | undefined): string {
  const date = convertTimestampToDate(timestamp);
  if (!date) {
    return '-';
  }

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

/**
 * Format protobuf Timestamp to date only string
 * Format: YYYY-MM-DD
 */
export function formatTimestampDate(timestamp: Timestamp | undefined): string {
  const date = convertTimestampToDate(timestamp);
  if (!date) {
    return '-';
  }

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');

  return `${year}-${month}-${day}`;
}

/**
 * Format protobuf Timestamp to relative time string
 * Example: "2小时前", "3天前"
 */
export function formatRelativeTime(timestamp: Timestamp | undefined): string {
  const date = convertTimestampToDate(timestamp);
  if (!date) {
    return '-';
  }

  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffSecs = Math.floor(diffMs / 1000);
  const diffMins = Math.floor(diffSecs / 60);
  const diffHours = Math.floor(diffMins / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffSecs < 60) {
    return '刚刚';
  } else if (diffMins < 60) {
    return `${diffMins}分钟前`;
  } else if (diffHours < 24) {
    return `${diffHours}小时前`;
  } else if (diffDays < 30) {
    return `${diffDays}天前`;
  } else {
    return formatTimestampDate(timestamp);
  }
}

/**
 * Check if timestamp is expired
 */
export function isTimestampExpired(timestamp: Timestamp | undefined): boolean {
  const date = convertTimestampToDate(timestamp);
  if (!date) {
    return true;
  }
  return date.getTime() < Date.now();
}

/**
 * Get remaining time until timestamp expires
 * Returns milliseconds, or 0 if already expired
 */
export function getRemainingTime(timestamp: Timestamp | undefined): number {
  const date = convertTimestampToDate(timestamp);
  if (!date) {
    return 0;
  }

  const remaining = date.getTime() - Date.now();
  return remaining > 0 ? remaining : 0;
}
