;; Wrapper functions around the pycontext/rbcontext command line utilities:
;; https://github.com/danielhones/context


(defvar context-keybinding-prefix "C-c q ")  ; need to have a space at the end

(defvar context-script-map
  '((ruby-mode "rbcontext")
    (python-mode "pycontext")))

(defun context-script ()
  (cadr (assoc major-mode context-script-map)))

(defun current-line-number ()
    (line-number-at-pos (point)))

(defun show-line-context ()
  "Show context for the current line"
  (interactive)
  (save-excursion
    (let* ((look-for (number-to-string (current-line-number)))
           (cmd (combine-and-quote-strings (list (context-script) "-nl" look-for)))
           (context-buffer (get-buffer-create "*Context*")))
      (mark-whole-buffer)
      (if look-for (shell-command-on-region (region-beginning) (region-end) cmd context-buffer)))))

(defun show-regex-context (look-for)
  "Show context using the look-for argument as a regex"
  (interactive "sLook for: ")
  (save-excursion
    (let ((cmd (combine-and-quote-strings (list (context-script) "-ne" look-for)))
          (context-buffer (get-buffer-create "*Context*")))
      (mark-whole-buffer)
      (if look-for (shell-command-on-region (region-beginning) (region-end) cmd context-buffer)))))

(defun show-at-point-context ()
  "Show context for the symbol currently at point"
  (interactive)
  (save-excursion
    (let* ((look-for (concat "'\b" (thing-at-point 'symbol t) "\b'"))
           (cmd (combine-and-quote-strings (list (context-script) "-ne" look-for)))
           (context-buffer (get-buffer-create "*Context*")))
      (mark-whole-buffer)
      (if look-for (shell-command-on-region (region-beginning) (region-end) cmd context-buffer)
        (message "no symbol at point")))))

(defun set-context-keybindings ()
  (global-set-key (kbd (concat context-keybinding-prefix "l")) 'show-line-context)
  (global-set-key (kbd (concat context-keybinding-prefix "c")) 'show-regex-context)
  (global-set-key (kbd (concat context-keybinding-prefix "p")) 'show-at-point-context))

(set-context-keybindings)
