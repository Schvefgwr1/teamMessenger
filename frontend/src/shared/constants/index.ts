export {
  ROUTES,
  PUBLIC_ROUTES,
  GUEST_ONLY_ROUTES,
  ADMIN_ROUTES,
} from './routes';

export {
  KNOWN_PERMISSIONS,
  KNOWN_PERMISSIONS as PERMISSIONS, // Алиас для удобства
  isKnownPermission,
  type KnownPermission,
} from './permissions';
