(defvar context-keybinding-prefix "C-c q ")  ; need to have a space at the end

(defvar context-script "code-context") ; Name or path of the context executable

(defvar context-language-map
  '((go-mode "go")
    (js-mode "js")
    (js2-mode "js")
    (python-mode "py")
    (ruby-mode "rb")
    (yaml-mode "yml")))

(defun context-language ()
  (cadr (assoc major-mode context-language-map)))

(defun current-line-number ()
  (line-number-at-pos (point)))

(defun run-command-in-context-window (source-code command)
  (let ((context-buffer (get-buffer-create "*Context*"))
        (context-window (get-buffer-window "*Context*"))
        (mode major-mode))
    (display-buffer context-buffer)
    (with-current-buffer context-buffer
      (setq buffer-read-only nil)
      (erase-buffer)
      (insert (concat command "\n"))
      ;; This next line can provide syntax-highlighting, but because these aren't completed
      ;; source files, it doesn't always look right:
      (funcall mode)
      (start-process-shell-command "context" context-buffer command)
      (process-send-string "context" source-code)
      (process-send-eof "context")
      (stop-process "context")
      (setq buffer-read-only t))))

(defun show-line-context ()
  "Show context for the current line"
  (interactive)
  (save-excursion
    (let* ((look-for (number-to-string (current-line-number)))
           (cmd (combine-and-quote-strings (list context-script "-l" (context-language) "-n" look-for))))
      (cond (look-for (run-command-in-context-window (buffer-string) cmd))))))

(defun show-regex-context (look-for)
  "Show context using the look-for argument as a regex"
  (interactive "sLook for: ")
  (save-excursion
    (let ((cmd (combine-and-quote-strings (list context-script "-l" (context-language) "-n" "-e" look-for))))
          (cond (look-for (run-command-in-context-window (buffer-string) cmd))))))

(defun show-at-point-context ()
  "Show context for the symbol currently at point"
  (interactive)
  (save-excursion
    (let* ((look-for (thing-at-point 'symbol t))
           ;; TODO: Use this once regex match is re-implemented:
           ;; (look-for (concat "\\b" (thing-at-point 'symbol t) "\\b"))
           (cmd (combine-and-quote-strings (list context-script "-l" (context-language) "-n" "-e" look-for))))
      (cond (look-for
             (run-command-in-context-window (buffer-string) cmd))
            (t
             (message "no symbol at point"))))))

(defun set-context-keybindings ()
  (global-set-key (kbd (concat context-keybinding-prefix "l")) 'show-line-context)
  (global-set-key (kbd (concat context-keybinding-prefix "c")) 'show-regex-context)
  (global-set-key (kbd (concat context-keybinding-prefix "p")) 'show-at-point-context))

(set-context-keybindings)
