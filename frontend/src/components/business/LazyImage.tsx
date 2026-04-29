/**
 * Lazy Image Component
 * Provides lazy loading for images with placeholder and error handling
 */

import { useState, useRef, useEffect } from 'react'
import { Spin } from '@arco-design/web-react'

interface LazyImageProps {
  src: string
  alt: string
  placeholder?: string
  className?: string
  style?: React.CSSProperties
  onClick?: () => void
  onLoad?: () => void
  onError?: () => void
}

export function LazyImage({
  src,
  alt,
  placeholder = 'https://via.placeholder.com/150',
  className = '',
  style,
  onClick,
  onLoad,
  onError,
}: LazyImageProps) {
  const [isLoaded, setIsLoaded] = useState(false)
  const [hasError, setHasError] = useState(false)
  const imgRef = useRef<HTMLImageElement>(null)
  const [shouldLoad, setShouldLoad] = useState(false)

  // Use Intersection Observer to detect when image is in viewport
  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setShouldLoad(true)
          observer.disconnect()
        }
      },
      {
        rootMargin: '50px', // Start loading 50px before image enters viewport
      }
    )

    const imgElement = imgRef.current
    if (imgElement) {
      observer.observe(imgElement)
    }

    return () => observer.disconnect()
  }, [])

  const handleLoad = () => {
    setIsLoaded(true)
    onLoad?.()
  }

  const handleError = () => {
    setHasError(true)
    setIsLoaded(true) // Still mark as loaded to show placeholder
    onError?.()
  }

  return (
    <div
      ref={imgRef}
      className={`relative overflow-hidden ${className}`}
      style={style}
      onClick={onClick}
    >
      {/* Loading spinner */}
      {!isLoaded && shouldLoad && (
        <div className="absolute inset-0 flex items-center justify-center bg-gray-100">
          <Spin />
        </div>
      )}

      {/* Actual image */}
      {shouldLoad && !hasError && (
        <img
          src={src}
          alt={alt}
          onLoad={handleLoad}
          onError={handleError}
          className={`w-full h-full object-cover transition-opacity duration-300 ${
            isLoaded ? 'opacity-100' : 'opacity-0'
          }`}
          loading="lazy"
        />
      )}

      {/* Placeholder */}
      {(!shouldLoad || hasError) && (
        <div
          className="absolute inset-0 flex items-center justify-center bg-gray-100 text-gray-400"
          style={{ backgroundImage: hasError ? undefined : `url(${placeholder})` }}
        >
          {hasError && (
            <div className="text-center">
              <div className="text-xs">加载失败</div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

/**
 * Hook for lazy loading multiple images
 * Useful for galleries and lists
 */
export function useLazyImageLoader() {
  const [loadedImages, setLoadedImages] = useState<Set<string>>(new Set())
  const [failedImages, setFailedImages] = useState<Set<string>>(new Set())

  const markLoaded = (src: string) => {
    setLoadedImages(prev => new Set(prev).add(src))
  }

  const markFailed = (src: string) => {
    setFailedImages(prev => new Set(prev).add(src))
  }

  const isLoaded = (src: string) => loadedImages.has(src)
  const hasFailed = (src: string) => failedImages.has(src)

  return {
    markLoaded,
    markFailed,
    isLoaded,
    hasFailed,
    loadedCount: loadedImages.size,
    failedCount: failedImages.size,
  }
}
