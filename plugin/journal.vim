if exists('g:loaded_journal') || &compatible
  finish
endif
let g:loaded_journal = 1

let s:journal_binary = 'markdown-journal'
let s:err_no_binary = s:journal_binary . ' executable not found'

if !exists('g:journal_timeline_cmd')
  let g:journal_timeline_cmd = 'timeline --level=2'
endif

if !exists('g:journal_labels_cmd')
  let g:journal_labels_cmd = 'labels --level=2'
endif

if !exists('g:journal_timeline_heading')
  let g:journal_timeline_heading = '# Timeline'
endif

if !exists('g:journal_labels_heading')
  let g:journal_labels_heading = '# Labels'
endif

if !exists('g:journal_timeline_file')
  let g:journal_timeline_file = 'timeline.md'
endif

if !exists('g:journal_lables_file')
  let g:journal_lables_file = 'lables.md'
endif

function! s:delete_section(start_pattern, end_pattern)
  let line_index = line('^')

  while line_index < line('$') && getline(line_index) !~# a:start_pattern
    let line_index += 1
  endwhile
  let start_index = line_index
  let line_index += 1

  while line_index < line('$') && getline(line_index) !~# a:end_pattern
    let line_index += 1
  endwhile
  let end_index = line_index - 1

  " Don't delete blank lines preceding the end_pattern.
  while end_index > start_index && getline(end_index) =~# '^$'
    let end_index -= 1
  endwhile

  silent execute start_index.','.end_index.'delete _'

  return start_index - 1
endfunction

function! s:timeline_contents()
  return [g:journal_timeline_heading]
        \ + systemlist(s:journal_binary . ' ' . g:journal_timeline_cmd)
endfunction

function! s:labels_contents()
  return [g:journal_labels_heading]
        \ + systemlist(s:journal_binary . ' ' . g:journal_labels_cmd)
endfunction

function! s:open_file(file)
  if buffer_name('%') != a:file
    execute 'edit ' . a:file
  endif
endfunction

function! s:timeline()
  if !executable(s:journal_binary)
    echohl WarningMsg | echo s:err_no_binary | echohl None
    return
  endif

  if exists('g:journal_timeline_file')
    call s:open_file(g:journal_timeline_file)
  endif

  let pos = s:delete_section('^' . g:journal_timeline_heading, '^# ')
  call append(pos, s:timeline_contents())
endfunction

function! s:labels()
  if !executable(s:journal_binary)
    echohl WarningMsg | echo s:err_no_binary | echohl None
    return
  endif

  if exists('g:journal_labels_file')
    call s:open_file(g:journal_labels_file)
  endif

  let pos = s:delete_section('^' . g:journal_labels_heading, '^# ')
  call append(pos, s:labels_contents())
endfunction

function! s:today()
  call s:open_file(strftime('%Y-%m-%d.md'))
endfunction

command! -nargs=0 JournalTimeline call s:timeline()
command! -nargs=0 JournalLabels call s:labels()
command! -nargs=0 JournalToday call s:today()
