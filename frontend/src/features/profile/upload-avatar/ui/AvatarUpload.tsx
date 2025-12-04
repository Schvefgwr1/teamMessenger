import { useState, useRef } from 'react';
import { Camera, X } from 'lucide-react';
import { cn } from '@/shared/lib/cn';
import { Avatar, toast } from '@/shared/ui';
import { formatFileSize, MAX_FILE_SIZE } from '@/shared/lib/formatFileSize';
import type { File as ApiFile } from '@/shared/types';

interface AvatarUploadProps {
  /** Текущий файл аватара из API */
  currentFile?: ApiFile | null;
  /** Fallback для инициалов */
  fallback?: string;
  /** Callback при выборе нового файла */
  onFileSelect: (file: File | null) => void;
  /** Размер аватара */
  size?: 'md' | 'lg' | 'xl';
  /** Отключить редактирование */
  disabled?: boolean;
  /** CSS класс */
  className?: string;
}

const sizeClasses = {
  md: 'w-20 h-20',
  lg: 'w-28 h-28',
  xl: 'w-36 h-36',
};

/**
 * Компонент загрузки аватара с превью
 */
export function AvatarUpload({
  currentFile,
  fallback,
  onFileSelect,
  size = 'lg',
  disabled = false,
  className,
}: AvatarUploadProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Валидация размера
    if (file.size > MAX_FILE_SIZE) {
      toast.error(`Файл слишком большой. Максимум: ${formatFileSize(MAX_FILE_SIZE)}`);
      return;
    }

    // Валидация типа
    if (!file.type.startsWith('image/')) {
      toast.error('Пожалуйста, выберите изображение');
      return;
    }

    // Создаём превью
    const reader = new FileReader();
    reader.onloadend = () => {
      setPreviewUrl(reader.result as string);
    };
    reader.readAsDataURL(file);

    setSelectedFile(file);
    onFileSelect(file);
  };

  const handleRemove = () => {
    setPreviewUrl(null);
    setSelectedFile(null);
    onFileSelect(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleClick = () => {
    if (!disabled) {
      fileInputRef.current?.click();
    }
  };

  // Определяем что показывать: превью нового файла или текущий аватар
  const showPreview = previewUrl !== null;

  return (
    <div className={cn('relative inline-block', className)}>
      {/* Hidden file input */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileChange}
        className="hidden"
        disabled={disabled}
      />

      {/* Avatar container */}
      <div
        className={cn(
          'relative rounded-full overflow-hidden cursor-pointer group',
          sizeClasses[size],
          disabled && 'cursor-not-allowed opacity-60'
        )}
        onClick={handleClick}
      >
        {showPreview ? (
          <img
            src={previewUrl}
            alt="Preview"
            className="w-full h-full object-cover"
          />
        ) : (
          <Avatar
            file={currentFile}
            fallback={fallback}
            className="w-full h-full"
          />
        )}

        {/* Overlay on hover */}
        {!disabled && (
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
            <Camera className="text-white" size={size === 'md' ? 20 : 28} />
          </div>
        )}
      </div>

      {/* Remove button */}
      {showPreview && !disabled && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            handleRemove();
          }}
          className="absolute -top-1 -right-1 p-1.5 rounded-full bg-error text-white hover:bg-red-600 transition-colors shadow-lg"
          title="Удалить"
        >
          <X size={14} />
        </button>
      )}

      {/* Info text */}
      {selectedFile && (
        <p className="text-xs text-neutral-400 mt-2 text-center truncate max-w-full">
          {selectedFile.name}
        </p>
      )}
    </div>
  );
}

