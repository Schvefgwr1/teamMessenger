import { useState, useRef, useEffect, KeyboardEvent } from 'react';
import { Send, Paperclip, X } from 'lucide-react';
import { Button, toast } from '@/shared/ui';
import { cn } from '@/shared/lib/cn';
import { formatFileSize, MAX_FILE_SIZE } from '@/shared/lib/formatFileSize';

interface MessageInputProps {
  onSend: (content: string, files?: File[]) => void;
  isLoading?: boolean;
  disabled?: boolean;
}

/**
 * Поле ввода сообщения с поддержкой файлов
 */
export function MessageInput({ onSend, isLoading, disabled }: MessageInputProps) {
  const [content, setContent] = useState('');
  const [files, setFiles] = useState<File[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // Инициализация высоты textarea при монтировании
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = '40px';
    }
  }, []);

  const handleSubmit = () => {
    const trimmedContent = content.trim();
    if (!trimmedContent && files.length === 0) return;

    onSend(trimmedContent, files.length > 0 ? files : undefined);
    setContent('');
    setFiles([]);

    // Сбрасываем высоту textarea к исходному размеру
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      textareaRef.current.style.height = '40px';
      textareaRef.current.style.overflowY = 'hidden';
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    // Отправка по Enter (без Shift)
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit();
    }
  };

  const handleTextareaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setContent(e.target.value);

    // Auto-resize textarea - работает в обе стороны (расширение и сжатие)
    const textarea = e.target;
    const minHeight = 40; // Минимальная высота (исходный размер)
    const maxHeight = 150; // Максимальная высота
    
    // Сбрасываем высоту для корректного расчета scrollHeight
    // Убираем inline style, чтобы браузер мог пересчитать размеры
    textarea.style.height = 'auto';
    
    // Вычисляем нужную высоту на основе содержимого
    const scrollHeight = textarea.scrollHeight;
    
    // Вычисляем новую высоту: минимум 40px, максимум 150px
    const newHeight = Math.max(minHeight, Math.min(scrollHeight, maxHeight));
    
    // Устанавливаем вычисленную высоту
    textarea.style.height = `${newHeight}px`;
    
    // Управляем overflow - показываем скролл только если достигнут максимум
    textarea.style.overflowY = scrollHeight > maxHeight ? 'auto' : 'hidden';
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = Array.from(e.target.files || []);

    // Валидация файлов
    const validFiles = selectedFiles.filter((file) => {
      if (file.size > MAX_FILE_SIZE) {
        toast.error(`${file.name} слишком большой (макс. ${formatFileSize(MAX_FILE_SIZE)})`);
        return false;
      }
      return true;
    });

    // Ограничение количества файлов
    const maxFiles = 5;
    if (files.length + validFiles.length > maxFiles) {
      toast.error(`Максимум ${maxFiles} файлов`);
      setFiles([...files, ...validFiles].slice(0, maxFiles));
    } else {
      setFiles([...files, ...validFiles]);
    }

    // Сброс input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const removeFile = (index: number) => {
    setFiles(files.filter((_, i) => i !== index));
  };

  const canSend = (content.trim().length > 0 || files.length > 0) && !isLoading && !disabled;

  return (
    <div className="border-t border-neutral-800 p-4 bg-neutral-900/50">
      {/* Прикрепленные файлы */}
      {files.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-3">
          {files.map((file, index) => (
            <div
              key={index}
              className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-neutral-800 text-sm"
            >
              <Paperclip size={14} className="text-neutral-400" />
              <span className="text-neutral-200 truncate max-w-[150px]">
                {file.name}
              </span>
              <span className="text-neutral-500 text-xs">
                ({formatFileSize(file.size)})
              </span>
              <button
                type="button"
                onClick={() => removeFile(index)}
                className="text-neutral-400 hover:text-error transition-colors"
              >
                <X size={14} />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Input area */}
      <div className="flex items-stretch gap-2" style={{ height: '40px' }}>
        {/* File input */}
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={handleFileSelect}
          className="hidden"
        />

        {/* File attach button */}
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          disabled={disabled}
          className={cn(
            'w-10 flex items-center justify-center rounded-xl transition-colors flex-shrink-0',
            'text-neutral-400 hover:text-neutral-200 hover:bg-neutral-800',
            'bg-neutral-800 border border-neutral-700',
            'p-0 m-0',
            disabled && 'opacity-50 cursor-not-allowed'
          )}
          style={{ height: '40px' }}
        >
          <Paperclip size={18} className="block" style={{ display: 'block', margin: '0' }} />
        </button>

        {/* Text input */}
        <div className="flex-1 relative">
          <textarea
            ref={textareaRef}
            value={content}
            onChange={handleTextareaChange}
            onKeyDown={handleKeyDown}
            placeholder="Введите сообщение..."
            disabled={disabled}
            rows={1}
            className={cn(
              'w-full px-4 rounded-xl resize-none',
              'bg-neutral-800 border border-neutral-700',
              'text-neutral-100 placeholder:text-neutral-500',
              'focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent',
              'transition-all duration-200',
              'max-h-[150px]',
              disabled && 'opacity-50 cursor-not-allowed'
            )}
            style={{ 
              minHeight: '40px',
              paddingTop: '10px',
              paddingBottom: '10px',
              lineHeight: '20px',
              boxSizing: 'border-box',
              resize: 'none'
            }}
          />
        </div>

        {/* Send button */}
        <Button
          onClick={handleSubmit}
          disabled={!canSend}
          isLoading={isLoading}
          size="icon"
          className="w-10 rounded-xl flex-shrink-0 !p-0"
          style={{ height: '40px' }}
        >
          <Send size={18} className="block" style={{ display: 'block', margin: '0' }} />
        </Button>
      </div>

      {/* Hint */}
      <p className="text-xs text-neutral-600 mt-2">
        Enter для отправки, Shift+Enter для новой строки
      </p>
    </div>
  );
}

