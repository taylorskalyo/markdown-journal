if exists('g:loaded_journal') || &compatible
  finish
endif
let g:loaded_journal = 1

let s:journal_binary = 'markdown-journal'
let s:err_no_binary = s:journal_binary . ' executable not found'

if !exists('g:journal_timeline_cmd')
  let g:journal_timeline_cmd = 'timeline --level=2'
endif

if !exists('g:journal_label_cmd')
  let g:journal_label_cmd = 'labels'
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

if !exists('g:journal_labels_file')
  let g:journal_labels_file = 'labels.md'
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

function! s:label_contents(label)
  let filter_arg = '--filter ' . a:label
  let contents = systemlist(s:journal_binary . ' ' . g:journal_label_cmd . ' ' . filter_arg)[1:-1]
  let contents[0] = substitute(contents[0], '^#* *.', '\U&', '')
  return contents
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

  let old_pos = line('.')
  let pos = s:delete_section('^' . g:journal_timeline_heading, '^# ')
  call append(pos, s:timeline_contents())
  execute 'normal ' . old_pos . 'G'
endfunction

function! s:labels()
  if !executable(s:journal_binary)
    echohl WarningMsg | echo s:err_no_binary | echohl None
    return
  endif

  if exists('g:journal_labels_file')
    call s:open_file(g:journal_labels_file)
  endif

  let old_pos = line('.')
  let pos = s:delete_section('^' . g:journal_labels_heading, '^# ')
  call append(pos, s:labels_contents())
  execute 'normal ' . old_pos . 'G'
endfunction

function! s:label(label)
  if empty(a:label)
    return
  endif

  if !executable(s:journal_binary)
    echohl WarningMsg | echo s:err_no_binary | echohl None
    return
  endif

  call s:open_file(a:label . '.md')

  let old_pos = line('.')
  let pos = s:delete_section('^' . a:label, '^# ')
  call append(pos, s:label_contents(a:label))
  execute 'normal ' . old_pos . 'G'
endfunction

function! s:today()
  call s:open_file(strftime('%Y-%m-%d.md'))
endfunction

function! s:card()
  let last_card = glob('*-*-*-*.md', 0, 1)[-1]
  let card_num = substitute(last_card, '\d\{4\}-\d\{2\}-\d\{2\}-\(\d\+\)\.md', '\1', '')
  let card_num += 1
  call s:open_file(strftime('%Y-%m-%d-' . card_num . '.md'))
endfunction

command! -nargs=0 JournalTimeline call s:timeline()
command! -nargs=* JournalLabels call s:labels()
command! -nargs=* JournalLabel call s:label(<q-args>)
command! -nargs=0 JournalToday call s:today()
command! -nargs=0 JournalCard call s:card()
