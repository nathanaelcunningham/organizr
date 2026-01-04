import { useNotificationStore } from '../stores/useNotificationStore';
import type { NotificationType } from '../stores/useNotificationStore';

/**
 * Hook to easily show toast notifications
 */
export function useToast() {
  const addNotification = useNotificationStore(
    (state) => state.addNotification
  );

  return {
    success: (message: string, duration?: number) =>
      addNotification('success', message, duration),
    error: (message: string, duration?: number) =>
      addNotification('error', message, duration),
    info: (message: string, duration?: number) =>
      addNotification('info', message, duration),
    warning: (message: string, duration?: number) =>
      addNotification('warning', message, duration),
    show: (type: NotificationType, message: string, duration?: number) =>
      addNotification(type, message, duration),
  };
}
