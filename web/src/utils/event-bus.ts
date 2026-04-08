/**
 * Event bus for module-to-module communication.
 * Implements a simple publish/subscribe pattern to reduce coupling.
 */

type EventHandler<T = unknown> = (data: T) => void;

class EventBus {
  private events: Record<string, EventHandler<unknown>[]> = {};

  /**
   * Subscribe to an event.
   * @param event Event name
   * @param handler Event handler
   * @returns Unsubscribe function
   */
  subscribe<T>(event: string, handler: EventHandler<T>) {
    if (!this.events[event]) {
      this.events[event] = [];
    }
    this.events[event].push(handler as EventHandler<unknown>);
    
    // Return an unsubscribe function.
    return () => this.unsubscribe(event, handler);
  }

  /**
   * Unsubscribe from an event.
   * @param event Event name
   * @param handler Event handler
   */
  unsubscribe<T>(event: string, handler: EventHandler<T>) {
    if (!this.events[event]) return;
    this.events[event] = this.events[event].filter(h => h !== (handler as EventHandler<unknown>));
  }

  /**
   * Publish an event.
   * @param event Event name
   * @param data Event payload
   */
  publish<T>(event: string, data: T) {
    if (!this.events[event]) return;
    this.events[event].forEach(handler => handler(data));
  }

  /**
   * Clear all handlers for a given event.
   * @param event Event name
   */
  clear(event: string) {
    if (this.events[event]) {
      delete this.events[event];
    }
  }

  /**
   * Clear all event handlers.
   */
  clearAll() {
    this.events = {};
  }
}

// Export a singleton instance.
export const eventBus = new EventBus();

export default eventBus;
